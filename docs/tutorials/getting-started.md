# Getting Started with cfnfind

This tutorial will teach you the basics of using cfnfind to find CloudFormation stacks in your AWS account.

## What You'll Learn

By the end of this tutorial, you'll know how to:
- Search for stacks by name
- Use AWS profiles
- Search specific regions
- Interpret cfnfind output

## Prerequisites

- cfnfind installed (see [Installation](./installation.md))
- AWS credentials configured
- At least one CloudFormation stack in your AWS account

## Step 1: Your First Search

Let's search for all stacks in your AWS account. We'll search for stacks containing the word "test":

```bash
cfnfind test
```

**What happens**: cfnfind searches all AWS regions for stacks with "test" in their name.

### Understanding the Output

The output is tab-separated with three columns:

```
stack-name          region       status
my-test-stack       us-east-1    CREATE_COMPLETE
test-api-stack      eu-west-1    UPDATE_COMPLETE
```

- **Column 1**: Stack name
- **Column 2**: AWS region where the stack exists
- **Column 3**: Current stack status

## Step 2: Search Specific Regions

Searching all regions can be slow. Let's search only specific regions:

```bash
cfnfind --region us-east-1 --region us-west-2 test
```

**What happens**: cfnfind only searches the specified regions (us-east-1 and us-west-2).

**Tip**: Use `--region` multiple times to specify multiple regions.

## Step 3: Using AWS Profiles

If you have multiple AWS accounts configured as profiles, specify which one to use:

```bash
cfnfind --profile production api
```

**What happens**: cfnfind uses the credentials from the "production" profile in your `~/.aws/credentials` file.

### Default Profile

If you don't specify `--profile`, cfnfind uses the "default" profile:

```bash
cfnfind api
# Equivalent to:
cfnfind --profile default api
```

## Step 4: Understanding Pattern Matching

cfnfind uses case-insensitive partial matching:

```bash
# Finds: my-api-stack, API-Gateway-Stack, legacy-api
cfnfind api
```

The pattern matches anywhere in the stack name:

- ✅ Prefix: `api-stack`
- ✅ Suffix: `my-api`
- ✅ Middle: `my-api-stack`
- ✅ Case-insensitive: `API-Stack`, `Api-Stack`

## Step 5: Handling No Results

If no stacks match your search:

```bash
cfnfind nonexistent
```

Output:
```
No stacks found
```

This means no stacks contain "nonexistent" in their name across all searched regions.

## Common Use Cases

### Find all stacks in a specific region

```bash
cfnfind --region ca-central-1 ""
```

**Note**: Use an empty string `""` to match all stacks.

### Find production stacks

```bash
cfnfind prod
```

### Search development account

```bash
cfnfind --profile dev database
```

### Search multiple accounts

```bash
cfnfind --profile production api
cfnfind --profile staging api
cfnfind --profile development api
```

## Troubleshooting

### "No stacks found" but you know stacks exist

Check:
1. Are you using the correct AWS profile?
2. Are you searching the right regions?
3. Does your IAM user have `cloudformation:DescribeStacks` permission?

### Permission errors

Ensure your IAM user/role has the required permissions. See [IAM Permissions](../reference/iam-permissions.md).

## Next Steps

- Learn more about [CLI options](../reference/cli-reference.md)
- Explore [How-to guides](../how-to/) for specific tasks
- Understand [how cfnfind works](../explanation/architecture.md)
