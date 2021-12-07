package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "github.com/jmccance/kin",
	Short: "A friendly CLI for working with Amazon Kinesis",
}

func init() {
	rootCmd.PersistentFlags().StringP("profile", "p", "", "AWS Profile Name")
	rootCmd.PersistentFlags().StringP("region", "r", "", "AWS Region Name")
}

func Execute() error {
	return rootCmd.Execute()
}
