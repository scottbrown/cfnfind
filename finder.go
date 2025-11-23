package cfnfind

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

type StackFinder struct {
	profile       string
	regions       []string
	clientFactory ClientFactory
}

func NewStackFinder(profile string, regions []string) *StackFinder {
	return &StackFinder{
		profile:       profile,
		regions:       regions,
		clientFactory: DefaultClientFactory,
	}
}

func NewStackFinderWithClientFactory(profile string, regions []string, clientFactory ClientFactory) *StackFinder {
	return &StackFinder{
		profile:       profile,
		regions:       regions,
		clientFactory: clientFactory,
	}
}

func (sf *StackFinder) FindStacks(ctx context.Context, pattern string) ([]Stack, error) {
	cfg, err := sf.loadAWSConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	regions, err := sf.resolveRegions(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve regions: %w", err)
	}

	return sf.searchStacksInRegions(ctx, cfg, regions, pattern)
}

func (sf *StackFinder) loadAWSConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(sf.profile))
}

func (sf *StackFinder) resolveRegions(ctx context.Context, cfg aws.Config) ([]string, error) {
	if len(sf.regions) > 0 {
		return sf.regions, nil
	}

	return getAllRegions(), nil
}

func (sf *StackFinder) searchStacksInRegions(ctx context.Context, cfg aws.Config, regions []string, pattern string) ([]Stack, error) {
	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results []Stack
		errs    []error
	)

	for _, region := range regions {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()

			stacks, err := sf.searchStacksInRegion(ctx, cfg, r, pattern)
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs = append(errs, fmt.Errorf("region %s: %w", r, err))
				return
			}

			results = append(results, stacks...)
		}(region)
	}

	wg.Wait()

	if len(errs) > 0 {
		return results, fmt.Errorf("errors occurred in some regions: %v", errs)
	}

	return results, nil
}

func (sf *StackFinder) searchStacksInRegion(ctx context.Context, cfg aws.Config, region, pattern string) ([]Stack, error) {
	regionalCfg := cfg.Copy()
	regionalCfg.Region = region

	client := sf.clientFactory(regionalCfg)

	var stacks []Stack
	var nextToken *string

	for {
		output, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to describe stacks: %w", err)
		}

		for _, stack := range output.Stacks {
			if sf.matchesPattern(aws.ToString(stack.StackName), pattern) {
				stacks = append(stacks, Stack{
					Name:   aws.ToString(stack.StackName),
					Region: region,
					Status: string(stack.StackStatus),
				})
			}
		}

		if output.NextToken == nil {
			break
		}
		nextToken = output.NextToken
	}

	return stacks, nil
}

func (sf *StackFinder) matchesPattern(stackName, pattern string) bool {
	return strings.Contains(strings.ToLower(stackName), strings.ToLower(pattern))
}

func getAllRegions() []string {
	return []string{
		"us-east-1", "us-east-2", "us-west-1", "us-west-2",
		"af-south-1", "ap-east-1", "ap-south-1", "ap-south-2",
		"ap-northeast-1", "ap-northeast-2", "ap-northeast-3",
		"ap-southeast-1", "ap-southeast-2", "ap-southeast-3", "ap-southeast-4",
		"ca-central-1", "ca-west-1",
		"eu-central-1", "eu-central-2", "eu-west-1", "eu-west-2", "eu-west-3",
		"eu-south-1", "eu-south-2", "eu-north-1",
		"il-central-1",
		"me-south-1", "me-central-1",
		"sa-east-1",
	}
}
