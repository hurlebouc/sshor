package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

var keepassPathFlag = flag.String("keepass", "", "path of the keepass vault")
var keepassIdFlag = flag.String("keepass-id", "", "entry in the keepass vault (/<PATH>/<OF>/<ENTRY> or /<PATH>/<OF>/<ENTRY>)")
var keepassPwdFlag = flag.String("keepass-pwd", "", "password of the keepass vault")
var loginFlag = flag.String("login", "", "SSH login")
var passwordFlag = flag.String("password", "", "SSH password")
var portFlag = flag.Uint("port", 0, "SSH port")

func printUsageAndExit(msg string) {
	fmt.Println(msg)
	fmt.Printf("sshor [options] HOST")
	flag.Usage()
	os.Exit(1)
}

func getFullHost() string {
	args := flag.Args()
	if len(args) == 0 {
		printUsageAndExit("Host is missing")
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
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.Username
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

func getPassword() string {
	if *passwordFlag != "" {
		return *passwordFlag
	}
	if *keepassPathFlag != "" {
		path := *keepassPathFlag
		pwd := *keepassPwdFlag
		id := *keepassIdFlag

		if id == "" {
			printUsageAndExit("Keepass ID access is empty")
		}
		return readKeepass(path, pwd, id, getLogin())
	}
	print("Password: ")
	var pwd string
	fmt.Scanln(&pwd)
	return pwd
}

func getAuthMethod() ssh.AuthMethod {
	pwd := getPassword()
	return ssh.Password(pwd)
}
