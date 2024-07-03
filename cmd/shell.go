/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/hurlebouc/sshor/config"
	"github.com/hurlebouc/sshor/ssh"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// shellCmd represents the shell command
var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Open a remote shell in SSH",
	Long:  `Open a remote shell in SSH`,
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		//fmt.Printf("debug cmd: %+v\n", cmd)
		//fmt.Printf("debug args: %+v\n", args)
		//fmt.Printf("debug toComplete: %s\n", toComplete)
		if len(args) == 0 {
			return findAllPossibleHosts(toComplete), cobra.ShellCompDirectiveDefault
		} else {
			return []string{}, cobra.ShellCompDirectiveDefault
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ReadConf()
		if err != nil {
			panic(fmt.Errorf("cannot read config: %w", err))
		}
		ssh.Shell(getLogin(args, config), getHost(args, config), getPort(args, config), getAuthMethod(args, config))
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// shellCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// shellCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func findAllPossibleHosts(toComplete string) []string {
	login, host, _ := splitFullHost(toComplete)

	config, err := config.ReadConf()
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
