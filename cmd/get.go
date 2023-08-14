/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
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

        client := netcli.NewDefaultHTTPClient()
        
        res, err := client.GetRequest(context.Background(), url)
        if err != nil {
            fmt.Printf("Error sending GET request: %v\n", err)
            return
        }

        if res.Error != nil {
            fmt.Printf("Error: %d - %s\n", res.Error.Status, res.Error.Message)
            return
        }

        fmt.Printf("Response:\n%s\n", res.Data)
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
