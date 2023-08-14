/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	netcli "github.com/kirre02/net-cli/pkg/net-cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// postCmd represents the post command
var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Runs a POST request over HTTP",
    RunE: func(cmd *cobra.Command, args []string) error {
        payloadData := viper.GetStringMapString("payload")

        payloadBytes, err := json.Marshal(payloadData)
        if err != nil {
            return err
        }

        var url string

        // Check if URL is provided as an argument
		if len(args) == 1 {
			url = args[0]
		} else {
			// Use URL from config
			url = viper.GetString("url")
		}

        response, err := netcli.PostRequest(url, payloadBytes)
        if err != nil {
            return err
        }

        fmt.Printf("Response Status: %d\n", *response.Status)
        fmt.Printf("Response Data: %s\n", response.Data)

        return nil
    },
}
func init() {
	rootCmd.AddCommand(postCmd)
}
