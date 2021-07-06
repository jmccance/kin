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
	ShardId                     *string
	PartitionKey                *string
	SequenceNumber              *string
	ApproximateArrivalTimestamp *time.Time
	EncryptionType              types.EncryptionType
	Data                        *interface{}
}

func init() {
	tailCmd.Flags().StringP("stream-name", "n", "", "Stream name (required)")
	tailCmd.Flags().StringP("shard", "s", "", "Shard id; if not specified, all shards will be tailed")
	tailCmd.MarkFlagRequired("stream-name")

	rootCmd.AddCommand(tailCmd)
}

var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Tail records from a Kinesis Data Stream",
	Long: `Continuously reads records from the target stream. Each record's payload will be
deserialized as JSON if possible; otherwise it will be returned as a base64-encoded string.`,
	Run: runTailCmd,
}

func runTailCmd(cmd *cobra.Command, args []string) {
	streamName, _ := cmd.Flags().GetString("stream-name")
	shardId, _ := cmd.Flags().GetString("shard")

	client, err := aws.GetKinesisClient()
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	records := make(chan *RecordOutput)

	if shardId != "" {
		go tailStreamShard(client, &streamName, &shardId, records)
	} else {
		shardIds, err := getShardIds(client, &streamName)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		for _, shardId := range shardIds {
			go tailStreamShard(client, &streamName, shardId, records)
		}
	}

	for record := range records {
		jsonBytes, _ := json.Marshal(record)
		fmt.Println(string(jsonBytes))
	}
}

func getShardIds(client *kinesis.Client, streamName *string) ([]*string, error) {
	output, err := client.ListShards(context.TODO(), &kinesis.ListShardsInput{
		StreamName: streamName,
	})
	if err != nil {
		return nil, err
	}

	var streamNames = []*string{}
	for _, shard := range output.Shards {
		streamNames = append(streamNames, shard.ShardId)
	}
	return streamNames, nil
}

func tailStreamShard(client *kinesis.Client, streamName, shardId *string, out chan *RecordOutput) error {
	shardIterator, err := getShardIterator(client, streamName, shardId)
	if err != nil {
		// FIXME What is the right way to handle the error? Right now I think we just totally ignore
		// it, which seems bad.
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	for {
		res, err := client.GetRecords(
			context.TODO(),
			&kinesis.GetRecordsInput{ShardIterator: shardIterator},
		)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		for _, record := range res.Records {
			var data interface{}

			err = json.Unmarshal(record.Data, &data)
			if err != nil {
				// If we can't decode it as JSON, fallback to base64-encoded binary
				// TODO Logging the error at debug-level could be informative
				data = record.Data
			}

			output := RecordOutput{
				ShardId:                     shardId,
				PartitionKey:                record.PartitionKey,
				SequenceNumber:              record.SequenceNumber,
				ApproximateArrivalTimestamp: record.ApproximateArrivalTimestamp,
				EncryptionType:              record.EncryptionType,
				Data:                        &data,
			}
			out <- &output
		}

		shardIterator = res.NextShardIterator
		if shardIterator == nil {
			break
		}

		time.Sleep(2 * time.Second)
	}

	return nil
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
