package ssh

import (
	"bytes"
	"context"
	"os"
	"os/user"

	"github.com/hurlebouc/sshor/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type PatternDetector struct {
	buffer   []byte
	position int // next byte to write
}

func (p *PatternDetector) addBytes(b []byte) {
	for i, v := range b {
		p.buffer[(i+p.position)%len(p.buffer)] = v
	}
	p.position = (p.position + len(b)) % len(p.buffer)
}

func (p PatternDetector) exportBytes() []byte {
	res := make([]byte, len(p.buffer))
	for i := range res {
		res[i] = p.buffer[(i+p.position)%len(p.buffer)]
	}
	return res
}

func Shell(hostConf config.Host, passwordFlag, keepassPwdFlag string) {
	keepassPwdMap := map[string]string{}
	if hostConf.GetKeepass() != nil && keepassPwdFlag != "" {
		keepassPwdMap[*hostConf.GetKeepass()] = keepassPwdFlag
	}

	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	ctx := context.WithValue(context.Background(), CURRENT_USER, currentUser.Username)

	conn, _ := newSshClient(ctx, hostConf, passwordFlag, keepassPwdMap)
	defer conn.Close()
	// Create a session
	session, err := conn.client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Set up terminal modes
	modes := ssh.TerminalModes{
		//ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		//ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
		panic(err)
	}

	//session.Stdin = os.Stdin
	input, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	//session.Stdout = os.Stdout
	output, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}
	cha := make(chan bool)
	session.Stderr = os.Stderr

	go func() {
		input.Write([]byte("su - testuser\n"))
		b := <-cha
		//time.Sleep(1 * time.Second)
		if b {
			input.Write([]byte("testuser\n"))
		}
		buffer := make([]byte, 5)
		for {
			n, err := os.Stdin.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				panic(err)
			}
			_, err = input.Write(buffer[0:n])
			if err != nil {
				panic(err)
			}
		}
	}()

	go func() {
		buffer := make([]byte, 5)
		p := PatternDetector{
			buffer: make([]byte, 100),
		}
		passed := false
		for {
			n, err := output.Read(buffer)
			if n == 0 {
				break
			}
			if err != nil {
				panic(err)
			}
			if !passed {
				p.addBytes(buffer[0:n])
				if bytes.Contains(p.exportBytes(), []byte("Mot de passe : ")) {
					cha <- true
					passed = true
				}
			}
			_, err = os.Stdout.Write(buffer[0:n])
			if err != nil {
				panic(err)
			}
		}
	}()

	// Start remote shell
	if err := session.Shell(); err != nil {
		panic(err)
	}
	session.Wait()
}
