package cmd_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type file struct {
	content []byte
}

type directoryLayout struct {
	files map[string]file
	dirs  map[string]directoryLayout
}

func equalsFile(a, b file) bool {
	if len(a.content) != len(b.content) {
		return false
	}
	for i, v := range a.content {
		if v != b.content[i] {
			return false
		}
	}
	return true
}

func equalDirs(a, b directoryLayout) bool {
	if len(a.files) != len(b.files) {
		return false
	}
	if len(a.dirs) != len(b.dirs) {
		return false
	}
	for k, v := range a.files {
		vv, present := b.files[k]
		if !present {
			return false
		}
		if !equalsFile(v, vv) {
			return false
		}
	}
	for k, v := range a.dirs {
		vv, present := b.dirs[k]
		if !present {
			return false
		}
		if !equalDirs(v, vv) {
			return false
		}
	}
	return true
}

func TestEqualsFiles(t *testing.T) {
	a := file{content: []byte{}}
	if !equalsFile(a, a) {
		t.Fatalf("file\n--> %+v\ninvalidate \"=\" reflexivity", a)
	}
	aa := file{content: []byte{}}
	if !equalsFile(a, aa) {
		t.Fatalf("file\n--> %+v\nand\n--> %+v\nshould be equal", a, aa)
	}
	b := file{content: []byte("coucou")}
	if equalsFile(a, b) {
		t.Fatalf("file\n--> %+v\nand\n--> %+v\nshould be differents", a, b)
	}
	bb := file{content: []byte("coucou")}
	if !equalsFile(b, bb) {
		t.Fatalf("file\n--> %+v\nand\n--> %+v\nshould be equal", b, bb)
	}
}

func TestEqualsDirs(t *testing.T) {
	a := directoryLayout{}
	if !equalDirs(a, a) {
		t.Fatalf("layout\n--> %+v\ninvalidate \"=\" reflexivity", a)
	}
	aa := directoryLayout{}
	if !equalDirs(a, aa) {
		t.Fatalf("layout\n--> %+v\nand\n--> %+v\nshould be equal", a, aa)
	}
	b := directoryLayout{
		files: map[string]file{
			"test1": {
				content: []byte("coucou"),
			},
			"test2": {
				content: []byte("plop"),
			},
			"test3": {},
		},
	}
	if equalDirs(a, b) {
		t.Fatalf("layout\n--> %+v\nand\n--> %+v\nshould be differents", a, b)
	}
	bb := directoryLayout{
		files: map[string]file{
			"test1": {
				content: []byte("coucou"),
			},
			"test2": {
				content: []byte("plop"),
			},
			"test3": {},
		},
	}
	if !equalDirs(b, bb) {
		t.Fatalf("layout\n--> %+v\nand\n--> %+v\nshould be equal", b, bb)
	}
}

func populateDirectory(layout directoryLayout, dirpath string) {
	for filename, fileLayout := range layout.files {
		f, err := os.Create(filepath.Join(dirpath, filename))
		if err != nil {
			panic(err)
		}
		_, err = f.Write(fileLayout.content)
		if err != nil {
			panic(err)
		}
	}
	for direname, dirLayout := range layout.dirs {
		subdir := filepath.Join(dirpath, direname)
		err := os.Mkdir(subdir, 0777)
		if err != nil {
			panic(err)
		}
		populateDirectory(dirLayout, subdir)
	}
}

func readDirectory(dirpath string) directoryLayout {
	entries, err := os.ReadDir(dirpath)
	if err != nil {
		panic(err)
	}
	dirs := map[string]directoryLayout{}
	files := map[string]file{}
	for _, entry := range entries {
		if entry.IsDir() {
			subLayout := readDirectory(filepath.Join(dirpath, entry.Name()))
			dirs[entry.Name()] = subLayout
		} else {
			bytes, err := os.ReadFile(filepath.Join(dirpath, entry.Name()))
			if err != nil {
				panic(err)
			}
			files[entry.Name()] = file{
				content: bytes,
			}
		}
	}
	return directoryLayout{
		files: files,
		dirs:  dirs,
	}
}

func initTempDir(dir directoryLayout) string {
	dirName, err := os.MkdirTemp(os.TempDir(), "sshor-test")
	if err != nil {
		panic(err)
	}
	populateDirectory(dir, dirName)
	return dirName
}

func startSftpServer(c chan struct{}, login, paswword string, port uint16, dir string) {
	debugStream := os.Stderr

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			fmt.Fprintf(debugStream, "Login: %s\n", c.User())
			if c.User() == login && string(pass) == paswword {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	rsaPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Cannot generate private key", err)
	}
	private, err := ssh.NewSignerFromKey(rsaPrivKey)
	if err != nil {
		log.Fatal("Cannot use rsa key for signing", err)
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:"+fmt.Sprintf("%d", port))
	if err != nil {
		log.Fatal("failed to listen for connection", err)
	}
	fmt.Printf("Listening on %v\n", listener.Addr())
	c <- struct{}{}

	nConn, err := listener.Accept()
	if err != nil {
		log.Fatal("failed to accept incoming connection", err)
	}

	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	_, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		log.Fatal("failed to handshake", err)
	}
	fmt.Fprintf(debugStream, "SSH server established\n")

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of an SFTP session, this is "subsystem"
		// with a payload string of "<length=4>sftp"
		fmt.Fprintf(debugStream, "Incoming channel: %s\n", newChannel.ChannelType())
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			fmt.Fprintf(debugStream, "Unknown channel type: %s\n", newChannel.ChannelType())
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Fatal("could not accept channel.", err)
		}
		fmt.Fprintf(debugStream, "Channel accepted\n")

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "subsystem" request.
		go func(in <-chan *ssh.Request) {
			for req := range in {
				fmt.Fprintf(debugStream, "Request: %v\n", req.Type)
				ok := false
				switch req.Type {
				case "subsystem":
					fmt.Fprintf(debugStream, "Subsystem: %s\n", req.Payload[4:])
					if string(req.Payload[4:]) == "sftp" {
						ok = true
					}
				}
				fmt.Fprintf(debugStream, " - accepted: %v\n", ok)
				req.Reply(ok, nil)
			}
		}(requests)

		serverOptions := []sftp.ServerOption{
			sftp.WithDebug(debugStream),
			sftp.WithServerWorkingDirectory(dir),
		}

		fmt.Fprintf(debugStream, "Read write server\n")

		server, err := sftp.NewServer(
			channel,
			serverOptions...,
		)
		if err != nil {
			log.Fatal(err)
		}
		if err := server.Serve(); err != nil {
			if err != io.EOF {
				log.Fatal("sftp server completed with error:", err)
			}
		}
		server.Close()
		log.Print("sftp client exited session.")
	}
}
