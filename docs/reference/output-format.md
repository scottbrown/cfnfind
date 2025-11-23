# Output Format

This document describes cfnfind's output format and how to interpret it.

## Format

cfnfind outputs results as **tab-separated values (TSV)** with three columns:

```
<stack-name>    <region>    <status>
```

## Columns

### Column 1: Stack Name

The full name of the CloudFormation stack.

**Example**: `my-api-stack-production`

### Column 2: Region

The AWS region where the stack exists.

**Example**: `us-east-1`

### Column 3: Status

The current status of the CloudFormation stack.

**Example**: `CREATE_COMPLETE`

## Stack Statuses

### Successful States

| Status | Description |
|--------|-------------|
| `CREATE_COMPLETE` | Stack was successfully created |
| `UPDATE_COMPLETE` | Stack was successfully updated |
| `IMPORT_COMPLETE` | Resources were successfully imported |

### In-Progress States

| Status | Description |
|--------|-------------|
| `CREATE_IN_PROGRESS` | Stack creation in progress |
| `UPDATE_IN_PROGRESS` | Stack update in progress |
| `DELETE_IN_PROGRESS` | Stack deletion in progress |
| `ROLLBACK_IN_PROGRESS` | Rollback in progress |
| `UPDATE_ROLLBACK_IN_PROGRESS` | Update rollback in progress |
| `REVIEW_IN_PROGRESS` | Stack is being reviewed |
| `IMPORT_IN_PROGRESS` | Import operation in progress |

### Failed States

| Status | Description |
|--------|-------------|
| `CREATE_FAILED` | Stack creation failed |
| `UPDATE_FAILED` | Stack update failed |
| `DELETE_FAILED` | Stack deletion failed |
| `ROLLBACK_FAILED` | Rollback failed |
| `UPDATE_ROLLBACK_FAILED` | Update rollback failed |
| `IMPORT_ROLLBACK_FAILED` | Import rollback failed |

### Rollback States

| Status | Description |
|--------|-------------|
| `ROLLBACK_COMPLETE` | Stack rolled back to previous state |
| `UPDATE_ROLLBACK_COMPLETE` | Update rolled back successfully |
| `IMPORT_ROLLBACK_COMPLETE` | Import rolled back successfully |

## Example Output

```
my-api-stack              us-east-1    CREATE_COMPLETE
my-api-stack              eu-west-1    CREATE_COMPLETE
production-database       us-west-2    UPDATE_COMPLETE
staging-infrastructure    ca-central-1 ROLLBACK_COMPLETE
dev-test-env             ap-southeast-1 DELETE_IN_PROGRESS
```

## No Results

When no stacks match the search pattern:

```
No stacks found
```

This is written to stdout (not stderr).

## Parsing Output

### Shell (awk)

Extract just stack names:

```bash
cfnfind api | awk '{print $1}'
```

Extract stacks in a specific region:

```bash
cfnfind api | awk '$2 == "us-east-1" {print $1}'
```

### Shell (cut)

Get only stack names:

```bash
cfnfind api | cut -f1
```

### Shell (grep)

Filter by status:

```bash
cfnfind api | grep "CREATE_COMPLETE"
```

### Python

```python
import subprocess

result = subprocess.run(
    ['cfnfind', 'api'],
    capture_output=True,
    text=True
)

for line in result.stdout.strip().split('\n'):
    if line and line != "No stacks found":
        name, region, status = line.split('\t')
        print(f"Stack: {name}, Region: {region}, Status: {status}")
```

### Shell Script

```bash
#!/bin/bash

cfnfind "$1" | while IFS=$'\t' read -r name region status; do
    if [ "$status" = "CREATE_COMPLETE" ]; then
        echo "Active stack: $name in $region"
    fi
done
```

## Output to File

### Save to file

```bash
cfnfind api > stacks.txt
```

### Save to CSV

While cfnfind outputs TSV, you can convert to CSV:

```bash
cfnfind api | tr '\t' ',' > stacks.csv
```

With headers:

```bash
echo "name,region,status" > stacks.csv
cfnfind api | tr '\t' ',' >> stacks.csv
```

## Sorting and Filtering

### Sort by stack name

```bash
cfnfind api | sort
```

### Sort by region

```bash
cfnfind api | sort -k2
```

### Filter by region

```bash
cfnfind api | grep -E "us-east-1|us-west-2"
```

### Count stacks

```bash
cfnfind api | grep -v "No stacks found" | wc -l
```

## Exit Status

cfnfind uses standard Unix exit codes:

- **0**: Success (stacks found or no stacks found)
- **1**: Error occurred

Check exit status in scripts:

```bash
if cfnfind api > /dev/null; then
    echo "Search completed successfully"
else
    echo "Error occurred" >&2
    exit 1
fi
```

## Future Output Formats

Currently, cfnfind only supports TSV output. Future versions may support:

- JSON output (`--format json`)
- CSV output (`--format csv`)
- Table output (`--format table`)

## See Also

- [CLI Reference](./cli-reference.md) - Complete CLI documentation
- [How-to: Export Results](../how-to/export-results.md) - Exporting search results
