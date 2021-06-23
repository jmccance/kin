package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	tailCmd.Flags().StringP("stream", "s", "", "Stream name (required)")
	tailCmd.MarkFlagRequired("stream")

	rootCmd.AddCommand(tailCmd)
}

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail records from a Kinesis Data Stream",

	Run: func(cmd *cobra.Command, args []string) {
		streamName, _ := cmd.Flags().GetString("stream")
		fmt.Printf("TBD, but would pull events from %s\n", streamName)
	},
}
