# cfnfind

CLI tool to find AWS CloudFormation stacks by full or partial name across one or multiple regions.

## Installation

### Using Go Install

```bash
go install github.com/scottbrown/cfnfind/cmd/cfnfind@latest
```

### Building from Source

```bash
git clone https://github.com/scottbrown/cfnfind.git
cd cfnfind
task build
# Binary will be in bin/cfnfind
```

## Prerequisites

- Go 1.21 or later
- AWS credentials configured (via `~/.aws/credentials` or environment variables)
- Appropriate IAM permissions to list CloudFormation stacks

## Usage

### Basic Usage

Find stacks by name pattern:

```bash
cfnfind my-stack
```

### With AWS Profile

Specify an AWS profile:

```bash
cfnfind --profile production my-stack
```

If `--profile` is not specified, the `default` profile is used.

### Searching Specific Regions

Search in specific regions:

```bash
cfnfind --region us-east-1 --region us-west-2 my-stack
```

If `--region` is not specified, all available AWS regions are searched.

### Examples

Search for all stacks containing "api":

```bash
cfnfind api
```

Search for stacks in a specific profile and region:

```bash
cfnfind --profile dev --region ca-central-1 database
```

Search across multiple regions:

```bash
cfnfind --region us-east-1 --region eu-west-1 --region ap-southeast-1 prod
```

## Output Format

The tool outputs matching stacks in tab-separated format:

```
stack-name    region          status
```

Example output:

```
my-api-stack          us-east-1    CREATE_COMPLETE
my-api-stack-v2       us-west-2    UPDATE_COMPLETE
```

## Development

### Building

```bash
task build
```

### Running Tests

```bash
task test
```

### Checking Coverage

```bash
task coverage
```

### Formatting Code

```bash
task fmt
```

### Running All Checks

```bash
task lint
```

## IAM Permissions

The tool requires the following IAM permission:

- `cloudformation:DescribeStacks`

Example IAM policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "cloudformation:DescribeStacks"
      ],
      "Resource": "*"
    }
  ]
}
```

## License

See [LICENSE](LICENSE) file for details.
