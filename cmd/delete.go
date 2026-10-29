package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kirre02/net-cli/internal/httpclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [url]",
	Short: "Runs a DELETE request over HTTP",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var url string

		if len(args) == 1 {
			url = args[0]
		} else {
			url = viper.GetString("url")
		}

		if url == "" {
			return fmt.Errorf("URL must be provided as an argument or configured")
		}

		client := httpclient.NewClient(0)

		res, err := client.DeleteRequest(context.Background(), url)
		if err != nil {
			return fmt.Errorf("DELETE request failed: %w", err)
		}

		responseJSON, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return fmt.Errorf("formatting response as JSON: %w", err)
		}

		fmt.Printf("Response:\n%s\n", responseJSON)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
