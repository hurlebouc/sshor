/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hurlebouc/sshor/config"
	"github.com/hurlebouc/sshor/ssh"
	"github.com/samber/lo"

	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"github.com/spf13/cobra"
	sshlib "golang.org/x/crypto/ssh"
)

const Version = "0.1.0"

var keepassPathFlag string
var keepassIdFlag string
var keepassPwdFlag string
var loginFlag string
var passwordFlag string
var portFlag uint16

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: Version,
	Use:     "sshor",
	Short:   "Tailored SSH",
}

func existFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func readConf() (*config.Config, error) {
	ctx := cuecontext.New()

	configdir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	path := "sshor.cue"
	if !existFile(path) {
		path, err = filepath.Abs(configdir + "/sshor/config.cue")
		if err != nil {
			panic(err)
		}
	}
	if !existFile(path) {
		return nil, nil
	}

	// Load the package "example" from the current directory.
	// We don't need to specify a config in this example.
	insts := load.Instances([]string{path}, nil)

	// The current directory just has one file without any build tags,
	// and that file belongs to the example package.
	// So we get a single instance as a result.
	value := ctx.BuildInstance(insts[0])

	if value.Err() != nil {
		return nil, value.Err()
	}

	jsonBytes, err := value.MarshalJSON()
	if err != nil {
		return nil, err
	}

	config := config.Config{}
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func findAllPossibleHosts(toComplete string) []string {
	login, host, _ := splitFullHost(toComplete)

	config, err := readConf()
	if err != nil {
		panic(err)
	}
	if config == nil {
		return []string{}
	}

	keys := make([]string, 0, len(config.Hosts))
	for k := range config.Hosts {
		keys = append(keys, k)
	}

	return lo.Map(lo.Filter(keys, func(item string, idx int) bool { return strings.HasPrefix(item, host) }), func(item string, idx int) string {
		if login == nil {
			return item
		} else {
			return *login + "@" + item
		}
	})
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

func getLogin(args []string, config *config.Config) string {
	if loginFlag != "" {
		return loginFlag
	}
	loginFromHost, host, _ := splitFullHost(getFullHost(args))
	if loginFromHost != nil {
		return *loginFromHost
	}

	loginFromConfig := config.GetHost(host).GetUser()
	if loginFromConfig != nil {
		return *loginFromConfig
	}

	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.Username

}

func getHost(args []string, config *config.Config) string {
	_, host, _ := splitFullHost(getFullHost(args))

	hostFromConfig := config.GetHost(host).GetHost()
	if hostFromConfig != nil {
		return *hostFromConfig
	}

	return host
}

func getPort(args []string, config *config.Config) uint16 {
	if portFlag != 0 {
		return uint16(portFlag)
	}
	_, host, portFromHost := splitFullHost(getFullHost(args))
	if portFromHost != nil {
		return *portFromHost
	}

	portFromConfig := config.GetHost(host).GetPort()
	if portFromConfig != nil {
		return *portFromConfig
	}

	return 22
}

func getPassword(args []string, config *config.Config) string {
	if passwordFlag != "" {
		return passwordFlag
	}

	_, host, _ := splitFullHost(getFullHost(args))

	var path string
	var pwd string
	var id string

	if keepassPathFlag != "" {
		path = keepassPathFlag
	} else {
		keepassFromConfig := config.GetHost(host).GetKeepass()
		if keepassFromConfig != nil {
			path = *keepassFromConfig
		}
	}

	if keepassPwdFlag != "" {
		pwd = keepassPwdFlag
	} else {
		pwdFromConfig := config.GetHost(host).GetKeepassPwd()
		if pwdFromConfig != nil {
			pwd = *pwdFromConfig
		}
	}

	if keepassIdFlag != "" {
		id = keepassIdFlag
	} else {
		idFromConfig := config.GetHost(host).GetKeepassId()
		if idFromConfig != nil {
			id = *idFromConfig
		}
	}

	if path != "" {
		if id == "" {
			panic("Keepass ID access is empty")
		}
		if pwd == "" {
			pwd = ssh.GetPassword("Keepass vault password: ")
		}
		return ssh.ReadKeepass(path, pwd, id, getLogin(args, config))
	}

	return ssh.GetPassword("Password: ")
}

func getAuthMethod(args []string, config *config.Config) sshlib.AuthMethod {
	pwd := getPassword(args, config)
	return sshlib.Password(pwd)
}
