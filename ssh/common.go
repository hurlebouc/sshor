package ssh

import (
	"context"
	"fmt"
	"syscall"

	"github.com/hurlebouc/sshor/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type Key struct{ v string }

var CURRENT_USER = Key{
	v: "CURRENT_USER",
}

func GetCurrentUser(ctx context.Context) string {
	return ctx.Value(CURRENT_USER).(string)
}

func readPassword(prompt string) string {
	print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	println("")
	if err != nil {
		panic(err)
	}
	return string(bytePassword)
}

func getPassword(user string, config config.Host, keepassPwdMap map[string]string) string {

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
		return ReadKeepass(*path, pwd, *id, user)
	}

	return readPassword(fmt.Sprintf("Password for %s@%s:%d: ", user, *config.GetHost(), config.GetPortOrDefault(22)))
}

func getAuthMethod(user string, config config.Host, keepassPwdMap map[string]string) ssh.AuthMethod {
	pwd := getPassword(user, config, keepassPwdMap)
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

func getUser(ctx context.Context, hostConfig config.Host) (string, context.Context) {
	if hostConfig.GetUser() != nil {
		user := *hostConfig.GetUser()
		newctx := context.WithValue(ctx, CURRENT_USER, user)
		return user, newctx
	}
	return GetCurrentUser(ctx), ctx
}

func newSshClientConfig(ctx context.Context, hostConfig config.Host, passwordFlag string, keepassPwdMap map[string]string) (*ssh.ClientConfig, context.Context) {
	user, newctx := getUser(ctx, hostConfig)
	var authMethod ssh.AuthMethod
	if passwordFlag != "" {
		authMethod = ssh.Password(passwordFlag)
	} else {
		authMethod = getAuthMethod(user, hostConfig, keepassPwdMap)
	}

	clientConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return clientConfig, newctx
}

func newSshClient(ctx context.Context, hostConfig config.Host, passwordFlag string, keepassPwdMap map[string]string) (SshClient, context.Context) {
	if hostConfig.GetJump() != nil {
		jumpHost := *hostConfig.GetJump()
		jumpClient, newctx := newSshClient(ctx, jumpHost, "", keepassPwdMap)
		conn, err := jumpClient.client.Dial("tcp", fmt.Sprintf("%s:%d", *hostConfig.GetHost(), hostConfig.GetPortOrDefault(22)))
		if err != nil {
			panic(err)
		}
		clientConfig, newctx := newSshClientConfig(newctx, hostConfig, passwordFlag, keepassPwdMap)
		ncc, chans, reqs, err := ssh.NewClientConn(conn, *hostConfig.GetHost(), clientConfig)
		if err != nil {
			panic(err)
		}
		sClient := ssh.NewClient(ncc, chans, reqs)
		return SshClient{
			client: sClient,
			jump:   &jumpClient,
		}, newctx
	}

	clientConfig, newctx := newSshClientConfig(ctx, hostConfig, passwordFlag, keepassPwdMap)
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", *hostConfig.GetHost(), hostConfig.GetPortOrDefault(22)), clientConfig)
	if err != nil {
		panic(err)
	}
	return SshClient{
		client: conn,
		jump:   nil,
	}, newctx
}
