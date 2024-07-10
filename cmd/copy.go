/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/hurlebouc/sshor/config"
	"github.com/hurlebouc/sshor/ssh/copy"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// sftpCmd represents the sftp command
var sftpCmd = &cobra.Command{
	Use:   "copy",
	Short: "copy files from/to remote",
	Long:  "copy files from/to remote",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		configGlobal := readConf()
		env := map[string]copy.Endpoint{}
		endpoints := lo.Map(args, func(item string, idx int) copy.Endpoint { return getEndpoint(configGlobal, env, parseArg(item)) })
		dst := endpoints[len(endpoints)-1]
		srcs := endpoints[1:]
		for _, src := range srcs {
			copy.Copy(getOptions(), src, dst)
		}
	},
}

func init() {
	rootCmd.AddCommand(sftpCmd)
}

type fichier struct {
	path   string
	remote *remoteArg
}

type remoteArg struct {
	host string
	user *string
}

func (remoteArg remoteArg) getUserAtHost() string {
	if remoteArg.user == nil {
		return remoteArg.host
	}
	return fmt.Sprintf("%s@%s", *remoteArg.user, remoteArg.host)
}

func parseArg(arg string) fichier {
	split := strings.SplitN(arg, ":", 2)
	if len(split) == 0 {
		panic(fmt.Sprintf("cannot parse %s as file path", arg))
	}
	if len(split) == 1 {
		return fichier{
			path:   arg,
			remote: nil,
		}
	}
	remote := parseRemoteArg(split[0])
	path := split[1]
	if path == "" {
		path = "."
	}
	return fichier{
		path:   path,
		remote: &remote,
	}
}

func parseRemoteArg(arg string) remoteArg {
	split := strings.SplitN(arg, "@", 2)
	if len(split) == 0 {
		panic(fmt.Sprintf("cannot parse %s as remote", arg))
	}
	if len(split) == 1 {
		return remoteArg{
			host: arg,
			user: nil,
		}
	}
	return remoteArg{
		host: split[1],
		user: &split[0],
	}
}

func getEndpoint(configGlobal *config.Config, env map[string]copy.Endpoint, fichier fichier) copy.Endpoint {
	if fichier.remote == nil {
		return copy.NewLocal(fichier.path)
	}
	endpoint, present := env[fichier.remote.getUserAtHost()]
	if present {
		return endpoint
	}
	hostConfig := getHostConfig(configGlobal, fichier.remote.getUserAtHost())
	endpoint = copy.NewRemote(hostConfig, passwordFlag, keepassPwdFlag, fichier.path)
	env[fichier.remote.getUserAtHost()] = endpoint
	return endpoint
}
