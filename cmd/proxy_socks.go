/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hurlebouc/sshor/ssh"
	"github.com/spf13/cobra"
)

// proxyForwardCmd represents the proxy command
var socksCmd = &cobra.Command{
	Use:   "proxy-socks",
	Short: "Open socks server on local host serving requests from remote host",
	Long:  "Open socks server on local host serving requests from remote host",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ssh.Socks(getHostConfig(readConf(), args[0]), getOptions(), passwordFlag, keepassPwdFlag, proxyOptions.listeningIP, proxyOptions.listeningPort)
	},
}

func init() {
	rootCmd.AddCommand(socksCmd)

	socksCmd.Flags().StringVar(&proxyOptions.listeningIP, LISTENING_IP, "127.0.0.1", "local listening IP")
	socksCmd.Flags().Uint16Var(&proxyOptions.listeningPort, LISTENING_PORT, 0, "local listening port (default use random)")
}
