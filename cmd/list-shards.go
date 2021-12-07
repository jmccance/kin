package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmccance/kin/pkg/aws"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/spf13/cobra"
)

func init() {
	listShardsCmd.Flags().StringP("stream-name", "n", "", "Stream name (required)")
	listShardsCmd.MarkFlagRequired("stream-name")

	rootCmd.AddCommand(listShardsCmd)
}

var listShardsCmd = &cobra.Command{
	Use:     "list-shards",
	Aliases: []string{"lss"},
	Short:   "List shards",

	Run: func(cmd *cobra.Command, args []string) {
		streamName, _ := cmd.Flags().GetString("stream-name")
		profile , _ := cmd.Flags().GetString("profile")
		region , _ := cmd.Flags().GetString("region")

		client, err := aws.GetKinesisClient(aws.WithProfile(profile), aws.WithRegion(region))
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		output, err := client.ListShards(context.TODO(), &kinesis.ListShardsInput{
			StreamName: &streamName,
		})
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		for _, shard := range output.Shards {
			jsonBytes, _ := json.Marshal(shard)
			fmt.Println(string(jsonBytes))
		}
	},
}
