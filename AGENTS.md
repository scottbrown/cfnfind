# AGENTS.md

This file provides guidance to AI coding assistants when working with code in this repository.

## Project Overview

cfnfind is a CLI tool for finding AWS CloudFormation stacks by name pattern across multiple regions. It performs concurrent searches to minimize latency when searching all AWS regions.

## Building and Testing

Uses Go Task for all build operations:

```bash
# Build the binary (outputs to .build/cfnfind)
task build

# Run all tests with race detection and coverage
task test

# View test coverage report
task coverage

# Run a single test
go test -v -run TestStackFinder_matchesPattern

# Run all checks (fmt, vet, test)
task lint

# Clean build artifacts
task clean
```

Test coverage must be maintained at ≥70%.

## Architecture

### Core Components

**Package Structure:**
- `cmd/cfnfind/` - CLI entry point using Cobra
- Root package (`cfnfind`) - Business logic and AWS integration
- `client.go` - AWS CloudFormation client interface abstraction
- `finder.go` - Core search logic with concurrent region searching
- `stack.go` - Stack data type

**Key Design Patterns:**

1. **Interface-based AWS clients** - `CloudFormationAPI` interface allows testing without AWS credentials. Use `NewStackFinderWithClientFactory()` in tests with mock implementations.

2. **Concurrent region search** - `searchStacksInRegions()` spawns a goroutine per region using `sync.WaitGroup`. Results are collected with mutex protection.

3. **Dependency injection** - `StackFinder` accepts a `ClientFactory` function to create AWS clients, enabling testability without third-party mocking libraries.

### Testing Strategy

- Unit tests use interface mocks (see `mockCloudFormationClient` in `finder_integration_test.go`)
- Never use third-party mocking libraries
- Test pagination scenarios explicitly (`TestStackFinder_searchStacksInRegion_WithPagination`)
- Test concurrent behaviour and error aggregation

### Pattern Matching

Uses case-insensitive partial matching via `strings.Contains(strings.ToLower(stackName), strings.ToLower(pattern))`. This is intentionally simple - not regex, not wildcards.

## Documentation Structure

Documentation follows the Diátaxis framework in `docs/`:
- `tutorials/` - Learning-oriented guides
- `how-to/` - Task-oriented recipes
- `reference/` - Information-oriented specifications
- `explanation/` - Understanding-oriented discussions

When adding documentation, place it in the appropriate Diátaxis category.

## AWS Considerations

- Requires only `cloudformation:DescribeStacks` IAM permission
- Default behaviour searches all AWS regions (see `getAllRegions()` in `finder.go`)
- Handles pagination via manual loop (not using AWS SDK paginators) to work with the interface abstraction
- Searches run concurrently but respect AWS API rate limits through SDK defaults

## Development Notes

- Binary name: `cfnfind` (package name uses same)
- Module path: `github.com/scottbrown/cfnfind`
- Uses AWS SDK v2 for Go
- CLI built with `github.com/spf13/cobra`
- All region searching is concurrent by default for performance
