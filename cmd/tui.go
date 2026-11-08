package cmd

import (
	"github.com/kirre02/net-cli/internal/storage"
	"github.com/kirre02/net-cli/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive terminal interface",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewSQLite()
		if err != nil {
			return err
		}
		defer store.Close()

		return tui.Run(store)
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
