package cmd

import (
	"context"
	"kin/pkg/aws"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listStreamsCmd)
}

var listStreamsCmd = &cobra.Command{
	Use:     "list-streams",
	Aliases: []string{"ls"},
	Short:   "List Kinesis streams",

	Run: func(cmd *cobra.Command, args []string) {
		client, err := aws.GetKinesisClient()
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}

		output, err := client.ListStreams(context.TODO(), &kinesis.ListStreamsInput{})
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}

		for _, name := range output.StreamNames {
			cmd.Println(name)
		}
	},
}
