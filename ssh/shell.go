package ssh

import (
	"context"
	"os"
	"os/user"

	"github.com/hurlebouc/sshor/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

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

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Start remote shell
	if err := session.Shell(); err != nil {
		panic(err)
	}
	session.Wait()
}
