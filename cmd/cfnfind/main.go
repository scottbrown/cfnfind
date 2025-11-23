package main

import (
	"context"
	"fmt"
	"os"

	"github.com/scottbrown/cfnfind"
	"github.com/spf13/cobra"
)

var (
	profile string
	regions []string
)

var rootCmd = &cobra.Command{
	Use:   "cfnfind [stack-name-pattern]",
	Short: "Find AWS CloudFormation stacks by full or partial name",
	Long:  "CLI tool to find AWS CloudFormation stacks by full or partial name across one or multiple regions",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func init() {
	rootCmd.Flags().StringVar(&profile, "profile", "default", "AWS profile to use")
	rootCmd.Flags().StringSliceVar(&regions, "region", []string{}, "AWS regions to search (defaults to all regions if not specified)")
}

func run(cmd *cobra.Command, args []string) error {
	stackPattern := args[0]
	ctx := context.Background()

	finder := cfnfind.NewStackFinder(profile, regions)
	stacks, err := finder.FindStacks(ctx, stackPattern)
	if err != nil {
		return fmt.Errorf("error finding stacks: %w", err)
	}

	if len(stacks) == 0 {
		fmt.Println("No stacks found")
		return nil
	}

	for _, stack := range stacks {
		fmt.Println(stack)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
