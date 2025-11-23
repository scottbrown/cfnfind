package cfnfind

import "testing"

func TestStack_String(t *testing.T) {
	tests := []struct {
		name     string
		stack    Stack
		expected string
	}{
		{
			name: "basic stack",
			stack: Stack{
				Name:   "my-stack",
				Region: "us-east-1",
				Status: "CREATE_COMPLETE",
			},
			expected: "my-stack\tus-east-1\tCREATE_COMPLETE",
		},
		{
			name: "stack with hyphens",
			stack: Stack{
				Name:   "my-app-stack-prod",
				Region: "ca-central-1",
				Status: "UPDATE_COMPLETE",
			},
			expected: "my-app-stack-prod\tca-central-1\tUPDATE_COMPLETE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.stack.String()
			if got != tt.expected {
				t.Errorf("Stack.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}
