package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

var keepassPathFlag = flag.String("keepass", "", "path of the keepass vault")
var keepassIdFlag = flag.String("keepass-id", "", "entry in the keepass vault (/<PATH>/<OF>/<ENTRY> or /<PATH>/<OF>/<ENTRY>)")
var keepassPwdFlag = flag.String("keepass-pwd", "", "password of the keepass vault")
var loginFlag = flag.String("login", "", "SSH login")
var passwordFlag = flag.String("password", "", "SSH password")
var portFlag = flag.Uint("port", 0, "SSH port")

func printUsageAndExit() {
	fmt.Printf("sshor [options] HOST")
	flag.Usage()
	os.Exit(1)
}

func getFullHost() string {
	args := flag.Args()
	if len(args) == 0 {
		printUsageAndExit()
	}
	return args[0]
}

func splitFullHost(fullHost string) (*string, string, *uint16) {
	splits := strings.SplitN(fullHost, "@", 2)
	var login *string
	var port *uint16
	var host string
	var hostWithPort string
	if len(splits) == 1 {
		hostWithPort = fullHost
		login = nil
	} else if len(splits) == 2 {
		hostWithPort = splits[1]
		login = &splits[0]
	} else {
		panic("unreachable")
	}
	portSplit := strings.SplitN(hostWithPort, ":", 2)
	if len(portSplit) == 1 {
		port = nil
		host = hostWithPort
	} else if len(portSplit) == 2 {
		portint, err := strconv.Atoi(portSplit[1])
		portu16 := uint16(portint)
		if err != nil {
			panic(err)
		}
		port = &portu16
		host = portSplit[0]
	} else {
		panic("unreachable")
	}
	return login, host, port
}

func getLogin() string {
	if *loginFlag != "" {
		return *loginFlag
	}
	loginFromHost, _, _ := splitFullHost(getFullHost())
	if loginFromHost != nil {
		return *loginFromHost
	}
	printUsageAndExit()
	panic("unreachable")
}

func getHost() string {
	_, host, _ := splitFullHost(getFullHost())
	return host
}

func getPort() uint16 {
	if *portFlag != 0 {
		return uint16(*portFlag)
	}
	_, _, portFromHost := splitFullHost(getFullHost())
	if portFromHost != nil {
		return *portFromHost
	}
	return 22
}

func getPassword() *string {
	if *passwordFlag != "" {
		return passwordFlag
	}
	return nil
}

func getAuthMethod() ssh.AuthMethod {
	pwd := getPassword()
	if pwd != nil {
		return ssh.Password(*pwd)
	} else {
		return ssh.KeyboardInteractive(nil) //todo
	}
}

func main() {

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
