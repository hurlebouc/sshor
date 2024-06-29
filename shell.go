package main

import (
	"flag"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func shell() {

	flag.Parse()

	// Create client config
	config := &ssh.ClientConfig{
		User: getLogin(),
		Auth: []ssh.AuthMethod{
			getAuthMethod(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", getHost()+":"+strconv.Itoa(int(getPort())), config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// Create a session
	session, err := conn.NewSession()
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

	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", height, width, modes); err != nil {
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
