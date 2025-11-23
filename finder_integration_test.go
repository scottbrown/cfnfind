package cfnfind

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type mockCloudFormationClient struct {
	describeStacksFunc func(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
}

func (m *mockCloudFormationClient) DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	if m.describeStacksFunc != nil {
		return m.describeStacksFunc(ctx, params, optFns...)
	}
	return &cloudformation.DescribeStacksOutput{}, nil
}

func TestStackFinder_searchStacksInRegion(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{Region: "us-east-1"}

	tests := []struct {
		name          string
		pattern       string
		mockStacks    []types.Stack
		expectedCount int
		expectError   bool
	}{
		{
			name:    "single matching stack",
			pattern: "api",
			mockStacks: []types.Stack{
				{
					StackName:   aws.String("my-api-stack"),
					StackStatus: types.StackStatusCreateComplete,
				},
				{
					StackName:   aws.String("database-stack"),
					StackStatus: types.StackStatusCreateComplete,
				},
			},
			expectedCount: 1,
		},
		{
			name:    "multiple matching stacks",
			pattern: "stack",
			mockStacks: []types.Stack{
				{
					StackName:   aws.String("api-stack"),
					StackStatus: types.StackStatusCreateComplete,
				},
				{
					StackName:   aws.String("database-stack"),
					StackStatus: types.StackStatusUpdateComplete,
				},
			},
			expectedCount: 2,
		},
		{
			name:          "no matching stacks",
			pattern:       "nonexistent",
			mockStacks:    []types.Stack{},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockCloudFormationClient{
				describeStacksFunc: func(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
					return &cloudformation.DescribeStacksOutput{
						Stacks: tt.mockStacks,
					}, nil
				},
			}

			clientFactory := func(cfg aws.Config) CloudFormationAPI {
				return mockClient
			}

			finder := NewStackFinderWithClientFactory("default", []string{}, clientFactory)
			stacks, err := finder.searchStacksInRegion(ctx, cfg, "us-east-1", tt.pattern)

			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(stacks) != tt.expectedCount {
				t.Errorf("expected %d stacks, got %d", tt.expectedCount, len(stacks))
			}
		})
	}
}

func TestStackFinder_searchStacksInRegion_WithPagination(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{Region: "us-east-1"}

	callCount := 0
	mockClient := &mockCloudFormationClient{
		describeStacksFunc: func(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
			callCount++
			if callCount == 1 {
				return &cloudformation.DescribeStacksOutput{
					Stacks: []types.Stack{
						{
							StackName:   aws.String("stack-1"),
							StackStatus: types.StackStatusCreateComplete,
						},
					},
					NextToken: aws.String("token1"),
				}, nil
			}
			return &cloudformation.DescribeStacksOutput{
				Stacks: []types.Stack{
					{
						StackName:   aws.String("stack-2"),
						StackStatus: types.StackStatusCreateComplete,
					},
				},
			}, nil
		},
	}

	clientFactory := func(cfg aws.Config) CloudFormationAPI {
		return mockClient
	}

	finder := NewStackFinderWithClientFactory("default", []string{}, clientFactory)
	stacks, err := finder.searchStacksInRegion(ctx, cfg, "us-east-1", "stack")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(stacks) != 2 {
		t.Errorf("expected 2 stacks from pagination, got %d", len(stacks))
	}

	if callCount != 2 {
		t.Errorf("expected 2 API calls for pagination, got %d", callCount)
	}
}

func TestStackFinder_searchStacksInRegion_Error(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{Region: "us-east-1"}

	mockClient := &mockCloudFormationClient{
		describeStacksFunc: func(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
			return nil, errors.New("API error")
		},
	}

	clientFactory := func(cfg aws.Config) CloudFormationAPI {
		return mockClient
	}

	finder := NewStackFinderWithClientFactory("default", []string{}, clientFactory)
	_, err := finder.searchStacksInRegion(ctx, cfg, "us-east-1", "stack")

	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestStackFinder_searchStacksInRegions(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{}

	mockClient := &mockCloudFormationClient{
		describeStacksFunc: func(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
			return &cloudformation.DescribeStacksOutput{
				Stacks: []types.Stack{
					{
						StackName:   aws.String("test-stack"),
						StackStatus: types.StackStatusCreateComplete,
					},
				},
			}, nil
		},
	}

	clientFactory := func(cfg aws.Config) CloudFormationAPI {
		return mockClient
	}

	finder := NewStackFinderWithClientFactory("default", []string{}, clientFactory)
	regions := []string{"us-east-1", "us-west-2"}
	stacks, err := finder.searchStacksInRegions(ctx, cfg, regions, "test")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(stacks) != 2 {
		t.Errorf("expected 2 stacks (one per region), got %d", len(stacks))
	}
}

func TestStackFinder_searchStacksInRegions_WithErrors(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{}

	mockClient := &mockCloudFormationClient{
		describeStacksFunc: func(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
			return nil, errors.New("region error")
		},
	}

	clientFactory := func(cfg aws.Config) CloudFormationAPI {
		return mockClient
	}

	finder := NewStackFinderWithClientFactory("default", []string{}, clientFactory)
	regions := []string{"us-east-1"}
	_, err := finder.searchStacksInRegions(ctx, cfg, regions, "test")

	if err == nil {
		t.Error("expected error but got nil")
	}
}
