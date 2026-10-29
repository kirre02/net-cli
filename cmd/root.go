package cmd

import (
	"fmt"
	"os"

	"github.com/kirre02/net-cli/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "net-cli",
	Short: "A simple CLI tool for HTTP calls",
	Run: func(cmd *cobra.Command, args []string) {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			fmt.Printf("version %s\n", constants.AppVersion)
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"config file",
	)
	rootCmd.PersistentFlags().BoolP(
		"version",
		"v",
		false,
		"Print version",
	)
	rootCmd.Flags().BoolP(
		"toggle",
		"t",
		false,
		"Help message for toggle",
	)
}

func initConfig() {
	if cfgFile != "" {
		// Use the explicitly provided config file.
		viper.SetConfigFile(cfgFile)
	} else {
		// Look for config.yaml in the current directory.
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		// Allow the CLI to run without a config file.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return
		}

		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		os.Exit(1)
	}
}
