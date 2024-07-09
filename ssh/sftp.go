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
