package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SERVER")
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "server",
	}

	rootCmd.AddCommand(
		runCmd(),
	)

	return rootCmd
}
