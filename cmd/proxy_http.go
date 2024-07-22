/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hurlebouc/sshor/ssh"
	"github.com/spf13/cobra"
)

// proxyForwardCmd represents the proxy command
var httpCmd = &cobra.Command{
	Use:   "proxy-http",
	Short: "Open HTTP proxy server on local host serving requests from remote host",
	Long:  "Open HTTP proxy server on local host serving requests from remote host",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return findAllPossibleHosts(toComplete), cobra.ShellCompDirectiveDefault
		} else {
			return []string{}, cobra.ShellCompDirectiveDefault
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ssh.Http(getHostConfig(readConf(), args[0]), getOptions(), passwordFlag, keepassPwdFlag, proxyOptions.listeningIP, proxyOptions.listeningPort)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)

	httpCmd.Flags().StringVar(&proxyOptions.listeningIP, LISTENING_IP, "127.0.0.1", "local listening IP")
	httpCmd.Flags().Uint16Var(&proxyOptions.listeningPort, LISTENING_PORT, 0, "local listening port (default use random)")
}
