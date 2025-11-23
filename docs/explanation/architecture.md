# Architecture

This document explains how cfnfind is designed and why certain architectural decisions were made.

## Overview

cfnfind is a simple, focused CLI tool with a clear purpose: find CloudFormation stacks by name across AWS regions. Its architecture reflects this simplicity.

## High-Level Architecture

```
┌─────────────┐
│   CLI       │  (Cobra command-line interface)
│  (cmd/)     │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Finder    │  (Business logic)
│  (finder.go)│
└──────┬──────┘
       │
       ├──────┐
       ▼      ▼
┌──────────┐ ┌──────────┐
│ AWS SDK  │ │ Goroutines│
│ Client   │ │ (Parallel)│
└──────────┘ └──────────┘
```

## Components

### 1. CLI Layer (`cmd/cfnfind/main.go`)

**Purpose**: Parse command-line arguments and coordinate the search

**Responsibilities**:
- Parse flags using Cobra
- Validate user input
- Invoke the Finder
- Format and display results

**Design decision**: We use Cobra because it's the de facto standard for Go CLI tools, providing:
- Automatic help generation
- Flag parsing
- Consistent UX with other CLI tools

### 2. Business Logic (`finder.go`)

**Purpose**: Orchestrate the CloudFormation stack search

**Key functions**:

```go
type StackFinder struct {
    profile       string
    regions       []string
    clientFactory ClientFactory
}

func (sf *StackFinder) FindStacks(ctx context.Context, pattern string) ([]Stack, error)
```

**Design decisions**:

1. **Dependency injection**: `ClientFactory` allows testing without AWS credentials
2. **Separation of concerns**: Logic separated from AWS API calls
3. **Immutability**: `StackFinder` is configured at creation, not modified during execution

### 3. Concurrent Search

**Why concurrent?**

Searching 20+ AWS regions sequentially would be slow. Searching them in parallel significantly reduces total time.

**How it works**:

```go
for _, region := range regions {
    wg.Add(1)
    go func(r string) {
        defer wg.Done()
        stacks, err := sf.searchStacksInRegion(ctx, cfg, r, pattern)
        // ... collect results
    }(region)
}
wg.Wait()
```

**Performance**:
- Sequential: ~20 seconds (1 second per region × 20 regions)
- Parallel: ~2-3 seconds (limited by slowest region)

**Trade-offs**:
- ✅ Much faster for multi-region searches
- ⚠️ Higher concurrent API request rate
- ⚠️ Need to handle concurrent writes with mutex

### 4. AWS SDK Integration

**Why AWS SDK v2?**

AWS SDK for Go v2 provides:
- Better performance
- Modular design (import only what you need)
- Context support (cancellation, timeouts)
- Active maintenance

**Client creation**:

```go
func DefaultClientFactory(cfg aws.Config) CloudFormationAPI {
    return cloudformation.NewFromConfig(cfg)
}
```

**Design decision**: Interface-based design allows mocking for tests without third-party mocking libraries.

## Data Flow

```
User Input
    │
    ▼
Parse flags (profile, regions, pattern)
    │
    ▼
Load AWS Config
    │
    ▼
Resolve regions (specified or all)
    │
    ▼
Spawn goroutine per region
    │
    ├─────┬─────┬─────┬─────┐
    ▼     ▼     ▼     ▼     ▼
  [us-east-1] [eu-west-1] [ap-southeast-1] ...
    │     │     │     │     │
    ▼     ▼     ▼     ▼     ▼
  DescribeStacks API calls
    │     │     │     │     │
    ▼     ▼     ▼     ▼     ▼
  Filter by pattern
    │     │     │     │     │
    └─────┴─────┴─────┴─────┘
              │
              ▼
      Collect results
              │
              ▼
       Sort & Display
```

## Design Principles

### 1. Simplicity

**Principle**: Do one thing well

cfnfind has a single purpose: find CloudFormation stacks. It doesn't:
- Modify stacks
- Display stack details
- Manage stack resources
- Replace the AWS CLI

**Benefit**: Simple code, easy to maintain, predictable behaviour

### 2. Performance

**Principle**: Be fast by default

- Concurrent searches across regions
- Minimal memory allocation
- Stream results (don't buffer everything)

### 3. Testability

**Principle**: Easy to test without AWS

**How**:
- Interfaces for AWS clients
- Dependency injection
- Pure functions where possible
- No global state

**Result**: 81.6% test coverage without AWS credentials

### 4. Idiomatic Go

**Principle**: Follow Go conventions

- Standard project layout (`cmd/`, root package)
- Error handling (return errors, don't panic)
- Interfaces (small, focused)
- No third-party mocking libraries
- Avoid primitive obsession (Stack type, not map)

## Pattern Matching

**Why partial, case-insensitive matching?**

```go
func (sf *StackFinder) matchesPattern(stackName, pattern string) bool {
    return strings.Contains(
        strings.ToLower(stackName),
        strings.ToLower(pattern),
    )
}
```

**Design decision**: Balance between simplicity and usefulness

Alternatives considered:
- ❌ Exact match only - too restrictive
- ❌ Regex - too complex for most use cases
- ❌ Wildcard patterns - adds complexity
- ✅ Partial, case-insensitive - simple and sufficient for 90% of use cases

**Future enhancement**: Add `--exact` flag for exact matching

## Error Handling

**Philosophy**: Fail fast, report clearly

```go
if err != nil {
    return fmt.Errorf("failed to describe stacks: %w", err)
}
```

**Benefits**:
- Clear error messages with context
- Error wrapping (`%w`) preserves original error
- Errors bubble up to CLI layer for consistent handling

## Region Selection

**Default**: Search all regions

**Why?** Users often don't know which region their stack is in. Searching all regions is the most useful default.

**Cost**: Minimal. `DescribeStacks` is free and fast.

**Alternative**: Let users specify regions when they know them (`--region`)

## Configuration

**Current**: No configuration file

**Why?** cfnfind is simple enough that configuration would add complexity without much benefit.

**Future**: If needed, could add:
- `~/.cfnfind.yml` for default regions/profile
- Project-level `.cfnfind.yml`

## Dependencies

**Minimal external dependencies**:
- Cobra (CLI framework)
- AWS SDK v2 (CloudFormation client)

**Why minimal?**
- Faster builds
- Fewer security vulnerabilities
- Easier maintenance
- Smaller binary

## Performance Characteristics

### Time Complexity

- Best case: O(r) where r = number of regions
- Worst case: O(r × s) where s = stacks per region

Parallel execution reduces wall-clock time to approximately O(s) for the slowest region.

### Space Complexity

O(n) where n = total matching stacks across all regions

Results are collected in memory. For thousands of stacks, memory usage remains minimal (each Stack is ~100 bytes).

## Security Considerations

### 1. Least Privilege

cfnfind requires only `cloudformation:DescribeStacks` permission.

### 2. Credential Handling

Uses AWS SDK's standard credential chain:
1. Environment variables
2. Shared credentials file
3. IAM role (EC2, ECS, etc.)

Never stores or logs credentials.

### 3. Read-Only

cfnfind cannot modify infrastructure. It's safe to use in production.

## Future Enhancements

Potential architectural changes:

1. **Output formats**: JSON, CSV, table
2. **Filtering**: By status, tags, creation date
3. **Caching**: Cache region results for repeated searches
4. **Streaming**: Stream results as they arrive
5. **Stack details**: Option to show more information

Each would require careful design to maintain simplicity.

## See Also

- [Concurrent Search Explanation](./concurrent-search.md)
- [Pattern Matching Explanation](./pattern-matching.md)
- [CLI Reference](../reference/cli-reference.md)
