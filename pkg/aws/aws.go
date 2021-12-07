package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func WithProfile(profile string) config.LoadOptionsFunc {
	return config.WithSharedConfigProfile(profile)
}

func WithRegion(region string) config.LoadOptionsFunc {
	return config.WithRegion(region)
}

func GetKinesisClient(cfgOpts ...func(*config.LoadOptions) error) (*kinesis.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), cfgOpts...)
	if err != nil {
		return nil, err
	}

	return kinesis.NewFromConfig(cfg), err
}
