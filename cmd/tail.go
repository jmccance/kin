package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"kin/pkg/aws"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
	"github.com/spf13/cobra"
)

func init() {
	tailCmd.Flags().StringP("stream-name", "n", "", "Stream name (required)")
	tailCmd.Flags().StringP("shard", "s", "0", "Shard id")
	tailCmd.MarkFlagRequired("stream-name")

	rootCmd.AddCommand(tailCmd)
}

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail records from a Kinesis Data Stream",

	Run: func(cmd *cobra.Command, args []string) {
		streamName, _ := cmd.Flags().GetString("stream-name")
		shardId, _ := cmd.Flags().GetString("shard")

		client, err := aws.GetKinesisClient()
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		shardIteratorOutput, err := client.GetShardIterator(
			context.TODO(),
			&kinesis.GetShardIteratorInput{
				ShardId:           &shardId,
				ShardIteratorType: types.ShardIteratorTypeTrimHorizon,
				StreamName:        &streamName,
			},
		)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		shardIterator := shardIteratorOutput.ShardIterator

		for {
			res, err := client.GetRecords(
				context.TODO(),
				&kinesis.GetRecordsInput{
					ShardIterator: shardIterator,
				},
			)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			for _, record := range res.Records {
				jsonBytes, _ := json.Marshal(record)
				fmt.Println(string(jsonBytes))
			}

			shardIterator = res.NextShardIterator
			if shardIterator == nil {
				break
			}

			time.Sleep(2 * time.Second)
		}
	},
}
