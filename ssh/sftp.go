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

	srcReader, err := os.Open(src)
	if err != nil {
		panic(err)
	}

	f.ReadFrom(srcReader)

	if _, err := f.Write([]byte("Hello world!")); err != nil {
		log.Fatal(err)
	}
	f.Close()

	// check it's there
	fi, err := client.Lstat("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fi)
}
