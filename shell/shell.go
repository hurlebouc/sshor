package shell

import (
	"os"
	"strconv"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func Shell(login, host string, port uint16, authMethod ssh.AuthMethod) {

	// Create client config
	config := &ssh.ClientConfig{
		User: login,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", host+":"+strconv.Itoa(int(port)), config)
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

func GetPassword(prompt string) string {
	print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	println("")
	if err != nil {
		panic(err)
	}
	return string(bytePassword)
}
