# CLI Reference

Complete command-line interface reference for cfnfind.

## Synopsis

```
cfnfind [flags] <stack-name-pattern>
```

## Description

cfnfind searches for AWS CloudFormation stacks by name pattern across one or more AWS regions. The search is case-insensitive and matches partial stack names.

## Arguments

### stack-name-pattern

The pattern to search for in stack names.

- **Type**: String
- **Required**: Yes
- **Pattern matching**: Case-insensitive partial match

**Examples**:
- `api` - Matches "my-api", "API-Gateway", "legacy-api-v2"
- `prod` - Matches "prod-stack", "production-api", "my-prod"
- `""` - Matches all stacks (empty string)

## Flags

### --profile

Specify which AWS profile to use from your credentials file.

- **Type**: String
- **Default**: `default`
- **Example**: `--profile production`

The profile must exist in:
- `~/.aws/credentials`
- `~/.aws/config`

### --region

Specify one or more AWS regions to search.

- **Type**: String (repeatable)
- **Default**: All AWS regions
- **Example**: `--region us-east-1 --region us-west-2`

**Supported regions**:
- us-east-1, us-east-2, us-west-1, us-west-2
- ca-central-1, ca-west-1
- eu-central-1, eu-central-2, eu-west-1, eu-west-2, eu-west-3
- eu-south-1, eu-south-2, eu-north-1
- ap-east-1, ap-south-1, ap-south-2
- ap-northeast-1, ap-northeast-2, ap-northeast-3
- ap-southeast-1, ap-southeast-2, ap-southeast-3, ap-southeast-4
- af-south-1
- il-central-1
- me-south-1, me-central-1
- sa-east-1

See [AWS Regions](./aws-regions.md) for the complete list.

### -h, --help

Display help information.

- **Type**: Boolean
- **Example**: `cfnfind --help`

## Output Format

cfnfind outputs tab-separated values with three columns:

```
<stack-name>    <region>    <status>
```

**Example**:
```
my-api-stack          us-east-1    CREATE_COMPLETE
prod-database         us-west-2    UPDATE_COMPLETE
test-infrastructure   eu-west-1    ROLLBACK_COMPLETE
```

See [Output Format](./output-format.md) for details.

## Exit Codes

- **0**: Success (stacks found or no stacks found)
- **1**: Error occurred (invalid arguments, AWS API errors, permission errors)

## Examples

### Basic search

Search for stacks containing "api" across all regions:

```bash
cfnfind api
```

### Search specific regions

Search only in us-east-1 and eu-west-1:

```bash
cfnfind --region us-east-1 --region eu-west-1 database
```

### Use specific AWS profile

Search using the "production" profile:

```bash
cfnfind --profile production backend
```

### Combined flags

Search for "api" in specific regions using a specific profile:

```bash
cfnfind --profile staging --region us-east-1 --region us-west-2 api
```

### Match all stacks

Find all stacks in a specific region:

```bash
cfnfind --region ca-central-1 ""
```

## Environment Variables

cfnfind respects standard AWS environment variables:

- `AWS_PROFILE` - Default profile (overridden by `--profile`)
- `AWS_REGION` - Default region (overridden by `--region`)
- `AWS_ACCESS_KEY_ID` - AWS access key
- `AWS_SECRET_ACCESS_KEY` - AWS secret key
- `AWS_SESSION_TOKEN` - Session token for temporary credentials

## See Also

- [IAM Permissions](./iam-permissions.md) - Required AWS permissions
- [Output Format](./output-format.md) - Understanding output
- [Getting Started](../tutorials/getting-started.md) - Tutorial for beginners
