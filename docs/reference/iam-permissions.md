# IAM Permissions

This document describes the AWS IAM permissions required to use cfnfind.

## Required Permissions

cfnfind requires the following IAM permission:

- `cloudformation:DescribeStacks`

This permission allows cfnfind to list and describe CloudFormation stacks in your AWS account.

## Minimal IAM Policy

Here's a minimal IAM policy that grants only the required permission:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "CfnFindDescribeStacks",
      "Effect": "Allow",
      "Action": [
        "cloudformation:DescribeStacks"
      ],
      "Resource": "*"
    }
  ]
}
```

### Why Resource is "*"

The `DescribeStacks` action is a list operation that doesn't support resource-level permissions. It must use `"Resource": "*"` to function properly.

## Creating an IAM User for cfnfind

### Step 1: Create the Policy

1. Open the AWS IAM Console
2. Navigate to **Policies** → **Create Policy**
3. Switch to the **JSON** tab
4. Paste the minimal policy above
5. Name it `CfnFindReadOnly`
6. Add description: "Allows cfnfind to describe CloudFormation stacks"
7. Click **Create Policy**

### Step 2: Create or Update User

**For a new user**:
1. Navigate to **Users** → **Add User**
2. Enter username (e.g., `cfnfind-cli`)
3. Select **Programmatic access**
4. Attach the `CfnFindReadOnly` policy
5. Complete the wizard and save credentials

**For an existing user**:
1. Navigate to **Users** → Select user
2. Click **Add permissions** → **Attach policies directly**
3. Search for and select `CfnFindReadOnly`
4. Click **Add permissions**

## Using IAM Roles

### For EC2 Instances

Attach an instance profile with the `CfnFindReadOnly` policy:

```bash
# Create role
aws iam create-role \
  --role-name CfnFindEC2Role \
  --assume-role-policy-document file://trust-policy.json

# Attach policy
aws iam attach-role-policy \
  --role-name CfnFindEC2Role \
  --policy-arn arn:aws:iam::ACCOUNT-ID:policy/CfnFindReadOnly

# Create instance profile
aws iam create-instance-profile \
  --instance-profile-name CfnFindEC2Role

# Add role to instance profile
aws iam add-role-to-instance-profile \
  --instance-profile-name CfnFindEC2Role \
  --role-name CfnFindEC2Role
```

**trust-policy.json**:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

### For ECS Tasks

Include the policy in your task execution role:

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

## Cross-Account Access

To search CloudFormation stacks in multiple AWS accounts:

### Step 1: Create Role in Target Account

In the account you want to search, create a role that can be assumed:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::SOURCE-ACCOUNT-ID:root"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

Attach the `CfnFindReadOnly` policy to this role.

### Step 2: Configure AWS Profile

Add the role to your `~/.aws/config`:

```ini
[profile target-account]
role_arn = arn:aws:iam::TARGET-ACCOUNT-ID:role/CfnFindRole
source_profile = default
```

### Step 3: Use the Profile

```bash
cfnfind --profile target-account api
```

## Read-Only vs Full Access

cfnfind only requires **read** permissions. It never:
- Creates stacks
- Updates stacks
- Deletes stacks
- Modifies stack resources

For production environments, grant only the minimal required permissions shown above.

## Troubleshooting Permission Errors

### "Access Denied" Error

```
Error: failed to describe stacks: AccessDenied: User is not authorized to perform: cloudformation:DescribeStacks
```

**Solutions**:
1. Verify the IAM policy is attached to your user/role
2. Check the policy has `cloudformation:DescribeStacks` permission
3. Ensure `Resource` is set to `"*"` (required for list operations)
4. Wait a few minutes for IAM changes to propagate

### "No Stacks Found" vs Permission Error

- **"No stacks found"**: Successful query, no matching stacks
- **Access Denied error**: Missing IAM permissions

If you see "No stacks found" but expect results, check:
1. Correct AWS profile (`--profile`)
2. Correct regions (`--region`)
3. Stack name pattern matches actual stack names

## Security Best Practices

1. **Principle of Least Privilege**: Grant only `cloudformation:DescribeStacks`, not full CloudFormation access
2. **Use IAM Roles**: Prefer IAM roles over long-term credentials when possible
3. **Audit Access**: Use CloudTrail to monitor `DescribeStacks` API calls
4. **Rotate Credentials**: Regularly rotate access keys for IAM users
5. **MFA**: Consider requiring MFA for sensitive accounts

## See Also

- [AWS IAM Documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/)
- [CloudFormation API Reference](https://docs.aws.amazon.com/AWSCloudFormation/latest/APIReference/)
- [Getting Started Tutorial](../tutorials/getting-started.md)
