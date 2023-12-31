/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
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

		payload, err := json.Marshal(payloadData)
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

		client := netcli.NewDefaultHTTPClient()
		// Send the POST request
		res, err := client.PostRequest(context.Background(), url, payload)
		if err != nil {
			return err
		}
		// Convert the response to JSON-formatted string with indentation
		responseJSON, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling response to JSON: %v\n", err)
			return err
		}

		// Print the JSON-formatted response
		fmt.Printf("Response:\n%s\n", responseJSON)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(postCmd)
}
