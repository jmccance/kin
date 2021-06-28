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

type RecordOutput struct {
	PartitionKey                *string
	SequenceNumber              *string
	ApproximateArrivalTimestamp *time.Time
	EncryptionType              types.EncryptionType
	Data                        *interface{}
}

func init() {
	tailCmd.Flags().StringP("stream-name", "n", "", "Stream name (required)")
	tailCmd.Flags().StringP("shard", "s", "0", "Shard id")
	tailCmd.MarkFlagRequired("stream-name")

	rootCmd.AddCommand(tailCmd)
}

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail records from a Kinesis Data Stream",
	Run:   runTailCmd,
}

func runTailCmd(cmd *cobra.Command, args []string) {
	streamName, _ := cmd.Flags().GetString("stream-name")
	shardId, _ := cmd.Flags().GetString("shard")

	client, err := aws.GetKinesisClient()
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	shardIterator, err := getShardIterator(client, &streamName, &shardId)
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	for {
		res, err := client.GetRecords(
			context.TODO(),
			&kinesis.GetRecordsInput{ShardIterator: shardIterator},
		)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		for _, record := range res.Records {
			var data interface{}

			err = json.Unmarshal(record.Data, &data)
			if err != nil {
				// If we can't decode it as JSON, fallback to base64-encoded binary
				// TODO Logging the error at debug-level could be informative
				data = record.Data
			}

			var output = RecordOutput{
				PartitionKey:                record.PartitionKey,
				SequenceNumber:              record.SequenceNumber,
				ApproximateArrivalTimestamp: record.ApproximateArrivalTimestamp,
				EncryptionType:              record.EncryptionType,
				Data:                        &data,
			}
			jsonBytes, _ := json.Marshal(output)
			fmt.Println(string(jsonBytes))
		}

		shardIterator = res.NextShardIterator
		if shardIterator == nil {
			break
		}

		time.Sleep(2 * time.Second)
	}
}

func getShardIterator(client *kinesis.Client, streamName *string, shardId *string) (*string, error) {
	shardIteratorOutput, err := client.GetShardIterator(
		context.TODO(),
		&kinesis.GetShardIteratorInput{
			ShardId:           shardId,
			ShardIteratorType: types.ShardIteratorTypeTrimHorizon,
			StreamName:        streamName,
		},
	)

	if err != nil {
		return nil, err
	} else {
		return shardIteratorOutput.ShardIterator, nil
	}
}
