/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	netcli "github.com/kirre02/net-cli/pkg/net-cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Run a GET request over HTTP",
	Run: func(cmd *cobra.Command, args []string) {
        var url string

        // Check if URL is provided as an argument
		if len(args) == 1 {
			url = args[0]
		} else {
			// Use URL from config
			url = viper.GetString("url")
		}

        response, err := netcli.GetRequest(url)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }

        if response.Error != nil {
            fmt.Printf("Error: %d - %s\n", response.Error.Status, response.Error.Message)
            return
        }

        fmt.Printf("Response:\n%s\n", response.Data)
    },
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
