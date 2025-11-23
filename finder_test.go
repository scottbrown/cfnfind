package cfnfind

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestNewStackFinder(t *testing.T) {
	profile := "test-profile"
	regions := []string{"us-east-1", "us-west-2"}

	finder := NewStackFinder(profile, regions)

	if finder.profile != profile {
		t.Errorf("NewStackFinder() profile = %q, want %q", finder.profile, profile)
	}

	if len(finder.regions) != len(regions) {
		t.Errorf("NewStackFinder() regions length = %d, want %d", len(finder.regions), len(regions))
	}

	for i, region := range regions {
		if finder.regions[i] != region {
			t.Errorf("NewStackFinder() regions[%d] = %q, want %q", i, finder.regions[i], region)
		}
	}
}

func TestStackFinder_matchesPattern(t *testing.T) {
	finder := &StackFinder{}

	tests := []struct {
		name      string
		stackName string
		pattern   string
		expected  bool
	}{
		{
			name:      "exact match",
			stackName: "my-stack",
			pattern:   "my-stack",
			expected:  true,
		},
		{
			name:      "partial match",
			stackName: "my-application-stack",
			pattern:   "application",
			expected:  true,
		},
		{
			name:      "case insensitive match",
			stackName: "MyStack",
			pattern:   "mystack",
			expected:  true,
		},
		{
			name:      "no match",
			stackName: "production-stack",
			pattern:   "development",
			expected:  false,
		},
		{
			name:      "prefix match",
			stackName: "prod-api-stack",
			pattern:   "prod",
			expected:  true,
		},
		{
			name:      "suffix match",
			stackName: "api-stack-v2",
			pattern:   "v2",
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := finder.matchesPattern(tt.stackName, tt.pattern)
			if got != tt.expected {
				t.Errorf("matchesPattern(%q, %q) = %v, want %v", tt.stackName, tt.pattern, got, tt.expected)
			}
		})
	}
}

func TestStackFinder_resolveRegions(t *testing.T) {
	ctx := context.Background()
	cfg := aws.Config{}

	tests := []struct {
		name             string
		specifiedRegions []string
		wantAllRegions   bool
	}{
		{
			name:             "with specified regions",
			specifiedRegions: []string{"us-east-1", "us-west-2"},
			wantAllRegions:   false,
		},
		{
			name:             "no specified regions",
			specifiedRegions: []string{},
			wantAllRegions:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := &StackFinder{regions: tt.specifiedRegions}
			regions, err := finder.resolveRegions(ctx, cfg)

			if err != nil {
				t.Fatalf("resolveRegions() error = %v", err)
			}

			if tt.wantAllRegions {
				allRegions := getAllRegions()
				if len(regions) != len(allRegions) {
					t.Errorf("resolveRegions() returned %d regions, want %d", len(regions), len(allRegions))
				}
			} else {
				if len(regions) != len(tt.specifiedRegions) {
					t.Errorf("resolveRegions() returned %d regions, want %d", len(regions), len(tt.specifiedRegions))
				}
				for i, region := range tt.specifiedRegions {
					if regions[i] != region {
						t.Errorf("resolveRegions() regions[%d] = %q, want %q", i, regions[i], region)
					}
				}
			}
		})
	}
}

func TestGetAllRegions(t *testing.T) {
	regions := getAllRegions()

	if len(regions) == 0 {
		t.Error("getAllRegions() returned empty slice")
	}

	expectedRegions := []string{
		"us-east-1", "us-west-2", "ca-central-1", "eu-west-1",
	}

	for _, expected := range expectedRegions {
		found := false
		for _, region := range regions {
			if region == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("getAllRegions() missing expected region %q", expected)
		}
	}
}
