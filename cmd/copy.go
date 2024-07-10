/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

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
		files := lo.Map(args, func(item string, idx int) fichier { return parseArg(item) })
		panic(files)
	},
}

func init() {
	rootCmd.AddCommand(sftpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sftpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sftpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type fichier struct {
	path   string
	remote *remoteArg
}

type remoteArg struct {
	host string
	user *string
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
	return fichier{
		path:   split[1],
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
