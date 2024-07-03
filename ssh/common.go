package ssh

import (
	"strconv"

	"github.com/hurlebouc/sshor/config"
	"golang.org/x/crypto/ssh"
)

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
