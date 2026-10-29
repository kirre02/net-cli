package cmd

import (
	"fmt"
	"time"

	"github.com/kirre02/net-cli/internal/httpclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Run a GET request over HTTP",
	RunE: func(cmd *cobra.Command, args []string) error {
		var url string

		// Check if URL is provided as an argument
		if len(args) == 1 {
			url = args[0]
		} else {
			// Use URL from config
			url = viper.GetString("url")
		}

		client := httpclient.NewClient(10 * time.Second)

		res, err := client.GetRequest(cmd.Context(), url)
		if err != nil {
			return fmt.Errorf("GET request failed: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Status: %d\n", res.StatusCode)

		_, err = cmd.OutOrStdout().Write(res.Body)
		if err != nil {
			return fmt.Errorf("writing response body: %w", err)
		}

		if len(res.Body) > 0 && res.Body[len(res.Body)-1] != '\n' {
			fmt.Println(cmd.OutOrStdout())
		}

		return nil
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
