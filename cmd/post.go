package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kirre02/net-cli/internal/httpclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var postCmd = &cobra.Command{
	Use:   "post [url]",
	Short: "Runs a POST request over HTTP",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var rawURL string

		// Use the URL argument if provided; otherwise use config.
		if len(args) == 1 {
			rawURL = args[0]
		} else {
			rawURL = viper.GetString("url")
		}

		if rawURL == "" {
			return fmt.Errorf(
				"URL is required: provide it as an argument or configure url",
			)
		}

		// Read the configured payload and encode it as JSON.
		payloadData := viper.GetStringMapString("payload")

		payload, err := json.Marshal(payloadData)
		if err != nil {
			return fmt.Errorf("encoding request payload: %w", err)
		}

		client := httpclient.NewClient(10 * time.Second)

		// Specify that the request body contains JSON.
		headers := http.Header{}
		headers.Set("Content-Type", "application/json")

		res, err := client.PostRequestWithHeaders(
			cmd.Context(),
			rawURL,
			payload,
			headers,
		)
		if err != nil {
			return fmt.Errorf("POST request failed: %w", err)
		}

		out := cmd.OutOrStdout()

		if _, err := fmt.Fprintf(out, "Status: %d\n", res.StatusCode); err != nil {
			return fmt.Errorf("writing response status: %w", err)
		}

		if _, err := out.Write(res.Body); err != nil {
			return fmt.Errorf("writing response body: %w", err)
		}

		if len(res.Body) > 0 && res.Body[len(res.Body)-1] != '\n' {
			if _, err := fmt.Fprintln(out); err != nil {
				return fmt.Errorf("writing response newline: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(postCmd)
}
