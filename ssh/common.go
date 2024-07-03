package ssh

import (
	"fmt"
	"strconv"
	"syscall"

	"github.com/hurlebouc/sshor/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func readPassword(prompt string) string {
	print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	println("")
	if err != nil {
		panic(err)
	}
	return string(bytePassword)
}

func getPassword(config config.Host, keepassPwdFlag string) string {

	path := config.Keepass
	pwd := keepassPwdFlag
	id := config.KeepassId

	if path != "" {
		if id == "" {
			panic("Keepass ID access is empty")
		}
		if pwd == "" {
			pwd = readPassword(fmt.Sprintf("Password for %s: ", path))
		}
		return ReadKeepass(path, pwd, id, config.User)
	}

	return readPassword(fmt.Sprintf("Password for %s@%s:%v: ", config.User, config.Host, config.Port))
}

func getAuthMethod(config config.Host, keepassPwdFlag string) ssh.AuthMethod {
	pwd := getPassword(config, keepassPwdFlag)
	return ssh.Password(pwd)
}

func newSshClient(login, host string, port uint16, authMethod ssh.AuthMethod, jumpHost *config.Host) *ssh.Client {
	// Create client config
	clientConfig := &ssh.ClientConfig{
		User: login,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", host+":"+strconv.Itoa(int(port)), clientConfig)
	if err != nil {
		panic(err)
	}
	return conn
}
