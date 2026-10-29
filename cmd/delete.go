package cmd

import (
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

		res, err := client.DeleteRequest(cmd.Context(), url)
		if err != nil {
			return fmt.Errorf("DELETE request failed: %w", err)
		}
		if err != nil {
			return fmt.Errorf("PUT request failed: %w", err)
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
	rootCmd.AddCommand(deleteCmd)
}
