package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func main() {

	var keepassPath = flag.String("keepass", "", "path of the keepass vault")
	var keepassId = flag.String("keepass-id", "", "entry in the keepass vault (/<PATH>/<OF>/<ENTRY> or /<PATH>/<OF>/<ENTRY>)")
	var keepassPwd = flag.String("keepass-pwd", "", "password of the keepass vault")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Printf("sshor [options] <USER>@<HOST>[:<PORT>]")
		flag.Usage()
		os.Exit(1)
	}
	println(*keepassPath)
	println(*keepassId)
	println(*keepassPwd)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Create client config
	config := &ssh.ClientConfig{
		User: "XXXXX",
		Auth: []ssh.AuthMethod{
			ssh.Password("WWWWWWW"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", "localhost:22", config)
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
