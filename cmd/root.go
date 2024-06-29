/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"

	"sshor/shell"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var keepassPathFlag string
var keepassIdFlag string
var keepassPwdFlag string
var loginFlag string
var passwordFlag string
var portFlag uint16

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sshor",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		shell.Shell(getLogin(args), getHost(args), getPort(args), getAuthMethod(args))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&keepassPathFlag, "keepass", "", "path of the keepass vault")
	rootCmd.PersistentFlags().StringVar(&keepassIdFlag, "keepass-id", "", "entry in the keepass vault (/<PATH>/<OF>/<ENTRY> or /<PATH>/<OF>/<ENTRY>)")
	rootCmd.PersistentFlags().StringVar(&keepassPwdFlag, "keepass-pwd", "", "password of the keepass vault")
	rootCmd.PersistentFlags().StringVar(&loginFlag, "login", "", "SSH login")
	rootCmd.PersistentFlags().StringVar(&passwordFlag, "password", "", "SSH password")
	rootCmd.PersistentFlags().Uint16Var(&portFlag, "port", 0, "SSH port")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getFullHost(args []string) string {
	if len(args) == 0 {
		panic("Host is missing")
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

func getLogin(args []string) string {
	if loginFlag != "" {
		return loginFlag
	}
	loginFromHost, _, _ := splitFullHost(getFullHost(args))
	if loginFromHost != nil {
		return *loginFromHost
	}
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.Username
}

func getHost(args []string) string {
	_, host, _ := splitFullHost(getFullHost(args))
	return host
}

func getPort(args []string) uint16 {
	if portFlag != 0 {
		return uint16(portFlag)
	}
	_, _, portFromHost := splitFullHost(getFullHost(args))
	if portFromHost != nil {
		return *portFromHost
	}
	return 22
}

func getPassword(args []string) string {
	if passwordFlag != "" {
		return passwordFlag
	}
	if keepassPathFlag != "" {
		path := keepassPathFlag
		pwd := keepassPwdFlag
		id := keepassIdFlag

		if id == "" {
			panic("Keepass ID access is empty")
		}
		return shell.ReadKeepass(path, pwd, id, getLogin(args))
	}
	print("Password: ")
	var pwd string
	fmt.Scanln(&pwd)
	return pwd
}

func getAuthMethod(args []string) ssh.AuthMethod {
	pwd := getPassword(args)
	return ssh.Password(pwd)
}
