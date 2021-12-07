package cmd

import (
	"context"
	"fmt"
	"github.com/jmccance/kin/pkg/aws"
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
		profile , _ := cmd.Flags().GetString("profile")
		region , _ := cmd.Flags().GetString("region")

		client, err := aws.GetKinesisClient(aws.WithProfile(profile), aws.WithRegion(region))
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
			fmt.Println(name)
		}
	},
}
