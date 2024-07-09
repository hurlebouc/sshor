package ssh

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Ne pas oublier le defer après !
func NewSftp(conn *ssh.Client) *sftp.Client {
	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	return client
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
