package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kin",
	Short: "A friendly CLI for working with Amazon Kinesis",
}

func Execute() error {
	return rootCmd.Execute()
}
