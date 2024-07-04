package ssh

import (
	"fmt"
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

func getPassword(config config.Host, keepassPwdMap map[string]string) string {

	path := config.GetKeepass()
	id := config.GetKeepassId()

	if path != nil {
		if id == nil {
			panic("Keepass ID access is empty")
		}
		pwd, present := keepassPwdMap[*path]
		if !present {
			pwd = readPassword(fmt.Sprintf("Password for %s: ", *path))
			keepassPwdMap[*path] = pwd
		}
		return ReadKeepass(*path, pwd, *id, config.User)
	}

	return readPassword(fmt.Sprintf("Password for %s@%s:%d: ", *config.GetUser(), *config.GetHost(), config.GetPortOrDefault(22)))
}

func getAuthMethod(config config.Host, keepassPwdMap map[string]string) ssh.AuthMethod {
	pwd := getPassword(config, keepassPwdMap)
	return ssh.Password(pwd)
}

type SshClient struct {
	client *ssh.Client
	jump   *SshClient
}

func (c SshClient) Close() {
	c.client.Close()
	if c.jump != nil {
		c.jump.Close()
	}
}

func newSshClientConfig(hostConfig config.Host, passwordFlag string, keepassPwdMap map[string]string) *ssh.ClientConfig {
	var authMethod ssh.AuthMethod
	if passwordFlag != "" {
		authMethod = ssh.Password(passwordFlag)
	} else {
		authMethod = getAuthMethod(hostConfig, keepassPwdMap)
	}

	clientConfig := &ssh.ClientConfig{
		User: *hostConfig.GetUser(),
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return clientConfig
}

func newSshClient(hostConfig config.Host, passwordFlag string, keepassPwdMap map[string]string) SshClient {
	if hostConfig.GetJump() != nil {
		jumpHost := *hostConfig.GetJump()
		jumpClient := newSshClient(jumpHost, "", keepassPwdMap)
		conn, err := jumpClient.client.Dial("tcp", fmt.Sprintf("%s:%d", *hostConfig.GetHost(), hostConfig.GetPortOrDefault(22)))
		if err != nil {
			panic(err)
		}
		clientConfig := newSshClientConfig(hostConfig, passwordFlag, keepassPwdMap)
		ncc, chans, reqs, err := ssh.NewClientConn(conn, *hostConfig.GetHost(), clientConfig)
		if err != nil {
			panic(err)
		}
		sClient := ssh.NewClient(ncc, chans, reqs)
		return SshClient{
			client: sClient,
			jump:   &jumpClient,
		}
	}

	clientConfig := newSshClientConfig(hostConfig, passwordFlag, keepassPwdMap)
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", *hostConfig.GetHost(), hostConfig.GetPortOrDefault(22)), clientConfig)
	if err != nil {
		panic(err)
	}
	return SshClient{
		client: conn,
		jump:   nil,
	}
}
