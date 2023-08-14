/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/kirre02/net-cli/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "net-cli",
	Short: "A simple CLI tool for http calls",
	Run: func(cmd *cobra.Command, args []string) {
        showVersion, _ := cmd.Flags().GetBool("version")
        if showVersion {
            fmt.Printf("version %s\n", constants.AppVersion)
            return
        }
    },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
    cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /net-cli.json)")
    rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version number") 
    // Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
        viper.SetConfigType("json")
    if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config file: %s\n", err)
			os.Exit(1)
		}
    }
}


