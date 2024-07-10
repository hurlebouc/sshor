package copy

import (
	"log"
	"os"
	"path/filepath"

	"github.com/hurlebouc/sshor/ssh"
	"github.com/pkg/sftp"
)

// Ne pas oublier le defer après !
func newSftp(conn ssh.SshClient) *sftp.Client {
	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(conn.Client)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func copyFile(options ssh.Options, src, dst Endpoint) {
	if options.Verbose {
		log.Printf("%s --> %s\n", src.url(), dst.url())
	}
	err := dst.fileSystem.mkdirAll(filepath.Dir(dst.path))
	if err != nil {
		panic(err)
	}

	f, err := dst.fileSystem.create(dst.path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	srcReader, err := src.fileSystem.open(src.path)
	if err != nil {
		panic(err)
	}
	defer srcReader.Close()

	_, err = f.ReadFrom(srcReader)
	if err != nil {
		panic(err)
	}
}

func copyDir(options ssh.Options, src, dst Endpoint) {
	err := dst.fileSystem.mkdirAll(dst.path)
	if err != nil {
		panic(err)
	}
	entries, err := src.fileSystem.readDir(src.path)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			copyDir(options, src.join(entry.Name()), dst.join(entry.Name()))
		} else {
			copyFile(options, src.join(entry.Name()), dst.join(entry.Name()))
		}
	}
}

func Copy(options ssh.Options, src, dst Endpoint) {
	if src.isDir() {
		copyDir(options, src, CompleteDstPath(src, dst))
	} else {
		copyFile(options, src, CompleteDstPath(src, dst))
	}
}

func CopyFileToRemote(client *sftp.Client, src, dst string) {

	err := client.MkdirAll(filepath.Dir(dst))
	if err != nil {
		panic(err)
	}

	f, err := client.Create(dst)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	srcReader, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcReader.Close()

	_, err = f.ReadFrom(srcReader)
	if err != nil {
		panic(err)
	}
}

func CopyFileToLocal(client *sftp.Client, src, dst string) {

	err := os.MkdirAll(filepath.Dir(dst), os.ModeDir)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	srcReader, err := client.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcReader.Close()

	_, err = f.ReadFrom(srcReader)
	if err != nil {
		panic(err)
	}
}

// CopyDirToRemote(client, /A/a, /B/b) creates /B/b and copies content of /A/a into /B/b
func CopyDirToRemote(client *sftp.Client, src, dst string) {
	err := client.MkdirAll(dst)
	if err != nil {
		panic(err)
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			CopyDirToRemote(client, filepath.Join(src, entry.Name()), dst+"/"+entry.Name())
		} else {
			CopyFileToRemote(client, filepath.Join(src, entry.Name()), dst+"/"+entry.Name())
		}
	}
}

// CopyDirToLocal(client, /A/a, /B/b) creates /B/b and copies content of /A/a into /B/b
func CopyDirToLocal(client *sftp.Client, src, dst string) {
	err := os.MkdirAll(dst, os.ModeDir)
	if err != nil {
		panic(err)
	}
	entries, err := client.ReadDir(src)
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			CopyDirToLocal(client, src+"/"+entry.Name(), filepath.Join(dst, entry.Name()))
		} else {
			CopyFileToLocal(client, src+"/"+entry.Name(), filepath.Join(dst, entry.Name()))
		}
	}
}
