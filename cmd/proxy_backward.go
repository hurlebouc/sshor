/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hurlebouc/sshor/ssh"
	"github.com/spf13/cobra"
)

// proxyBackwardCmd represents the proxy command
var proxyBackwardCmd = &cobra.Command{
	Use:   "proxy-backward",
	Short: "Backward connections from remote to local",
	Long:  `This command opens a listening port on remote host and backwards each request on this port to a destination accessible from the local host.`,
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return findAllPossibleHosts(toComplete), cobra.ShellCompDirectiveDefault
		} else {
			return []string{}, cobra.ShellCompDirectiveDefault
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		ssh.BackwardProxy(getHostConfig(readConf(), args[0]), getOptions(), passwordFlag, keepassPwdFlag, proxyOptions.listeningIP, proxyOptions.listeningPort, proxyOptions.destinationIP, proxyOptions.destinationPort)
	},
}

func init() {
	rootCmd.AddCommand(proxyBackwardCmd)

	proxyBackwardCmd.Flags().StringVar(&proxyOptions.listeningIP, LISTENING_IP, "127.0.0.1", "remote listening IP")
	proxyBackwardCmd.Flags().Uint16Var(&proxyOptions.listeningPort, LISTENING_PORT, 0, "remote listening port (default use random)")
	proxyBackwardCmd.Flags().StringVar(&proxyOptions.destinationIP, DESTINATION_ADDR, "127.0.0.1", "destination address accessible from local host")
	proxyBackwardCmd.Flags().Uint16Var(&proxyOptions.destinationPort, DESTINATION_PORT, 0, "destination port accessible from local host")
	proxyBackwardCmd.MarkFlagRequired(DESTINATION_PORT)
}
