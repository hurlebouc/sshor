/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hurlebouc/sshor/config"
	"github.com/hurlebouc/sshor/ssh"

	"github.com/spf13/cobra"
)

const Version = "1.0.1"

var keepassPathFlag string
var keepassIdFlag string
var keepassPwdFlag string
var loginFlag string
var passwordFlag string
var portFlag uint16
var verboseFlag bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: Version,
	Use:     "sshor",
	Short:   "Tailored SSH",
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
	rootCmd.PersistentFlags().StringVarP(&loginFlag, "login", "l", "", "SSH login")
	rootCmd.PersistentFlags().StringVarP(&passwordFlag, "password", "w", "", "SSH password")
	rootCmd.PersistentFlags().Uint16VarP(&portFlag, "port", "p", 0, "SSH port")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "verbosity")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

func getHost(arg string, config *config.Config) *string {
	_, host, _ := splitFullHost(arg)

	if config.GetHost(host) == nil {
		return &host
	}

	return config.GetHost(host).Host
}

func readConf() *config.Config {
	configGlobal, err := config.ReadConf()
	if err != nil {
		panic(fmt.Errorf("cannot read config: %w", err))
	}
	return configGlobal
}

func getHostConfig(configGlobal *config.Config, userAtHostPort string) config.Host {
	login, hostname, port := splitFullHost(userAtHostPort)
	hostConf := configGlobal.GetHost(hostname)
	if hostConf == nil {
		hostConf = &config.Host{}
	}
	if loginFlag != "" {
		hostConf.User = &loginFlag
	}
	if login != nil {
		hostConf.User = login
	}
	hostConf.Host = getHost(userAtHostPort, configGlobal)
	if portFlag != 0 {
		hostConf.Port = &portFlag
	}
	if port != nil {
		hostConf.Port = port
	}
	if keepassPathFlag != "" {
		hostConf.Keepass = &config.Keepass{
			Path: keepassPathFlag,
			Id:   keepassIdFlag,
		}
	}
	return *hostConf
}

func getOptions() ssh.Options {
	return ssh.Options{
		Verbose: verboseFlag,
	}
}
