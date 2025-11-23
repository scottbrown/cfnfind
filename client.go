package cfnfind

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

type CloudFormationAPI interface {
	DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
}

type ClientFactory func(cfg aws.Config) CloudFormationAPI

func DefaultClientFactory(cfg aws.Config) CloudFormationAPI {
	return cloudformation.NewFromConfig(cfg)
}
