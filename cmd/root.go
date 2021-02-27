package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go_api_skeleton",
	Short: "go_api_skeleton ",
	Long:  `Golang API skeleton`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(apiCmd)
}

// initConfig reads in ENV variables if set.
func initConfig() {
	viper.AutomaticEnv()
}
