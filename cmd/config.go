/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/hurlebouc/sshor/config"
	"github.com/nwidger/jsoncolor"
	"github.com/spf13/cobra"
)

// proxyBackwardCmd represents the proxy command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display JSON configuration",
	Long:  "Display JSON configuration",
	Run: func(cmd *cobra.Command, args []string) {
		jsonBytes, err := config.ReadConfJson()
		if err != nil {
			panic(err)
		}
		if colored {
			var raw json.RawMessage
			err = json.Unmarshal(jsonBytes, &raw)
			if err != nil {
				panic(err)
			}
			res, err := jsoncolor.MarshalIndent(raw, "", "\t")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(res))
		} else {
			fmt.Println(string(jsonBytes))
		}

	},
}

var colored bool

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVar(&colored, "format", false, "enable formatted and colored output")
}
