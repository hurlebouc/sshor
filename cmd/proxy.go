/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// proxyForwardCmd represents the proxy command
var proxyForwardCmd = &cobra.Command{
	Use:   "proxy-forward",
	Short: "Forward connections from local to remote",
	Long:  `This command opens a listening port on local host and forwards each request on this port to a destination accessible from the remote host.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("proxy called")
	},
}

var proxyForwardOptions struct {
	listeningIP     string
	listeningPort   uint16
	destinationIP   string
	destinationPort uint16
}

const LISTENING_IP = "listening-ip"
const LISTENING_PORT = "listening-port"
const DESTINATION_IP = "destination-ip"
const DESTINATION_PORT = "destination-port"

func init() {
	rootCmd.AddCommand(proxyForwardCmd)

	proxyForwardCmd.Flags().StringVar(&proxyForwardOptions.listeningIP, LISTENING_IP, "127.0.0.1", "local listening IP")
	proxyForwardCmd.Flags().Uint16Var(&proxyForwardOptions.listeningPort, LISTENING_PORT, 0, "local listening port (default use random)")
	proxyForwardCmd.Flags().StringVar(&proxyForwardOptions.destinationIP, DESTINATION_IP, "127.0.0.1", "destination IP accessible from remote host")
	proxyForwardCmd.Flags().Uint16Var(&proxyForwardOptions.destinationPort, DESTINATION_PORT, 0, "destination IP accessible from remote host")
	proxyForwardCmd.MarkFlagRequired(DESTINATION_PORT)
}
