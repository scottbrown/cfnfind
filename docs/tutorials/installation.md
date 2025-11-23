# Installation Tutorial

This tutorial will guide you through installing cfnfind on your system.

## Prerequisites

Before you begin, ensure you have:

- Go 1.21 or later installed
- Git installed (for building from source)
- AWS credentials configured

## Option 1: Install using Go

The easiest way to install cfnfind is using Go's install command:

```bash
go install github.com/scottbrown/cfnfind/cmd/cfnfind@latest
```

This will install the latest version of cfnfind to your `$GOPATH/bin` directory.

### Verify Installation

Check that cfnfind is installed correctly:

```bash
cfnfind --help
```

You should see the help message displaying available commands and flags.

## Option 2: Build from Source

If you prefer to build from source:

### Step 1: Clone the Repository

```bash
git clone https://github.com/scottbrown/cfnfind.git
cd cfnfind
```

### Step 2: Build the Binary

Using Go Task:

```bash
task build
```

Or using Go directly:

```bash
go build -o bin/cfnfind ./cmd/cfnfind
```

### Step 3: Move to Your PATH

```bash
sudo mv bin/cfnfind /usr/local/bin/
```

Or add the `bin/` directory to your PATH:

```bash
export PATH=$PATH:$(pwd)/bin
```

### Step 4: Verify Installation

```bash
cfnfind --help
```

## Configuring AWS Credentials

cfnfind uses AWS credentials from:

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. Shared credentials file (`~/.aws/credentials`)
3. AWS IAM role (when running on EC2, ECS, etc.)

### Setting up AWS Credentials

If you haven't configured AWS credentials:

```bash
aws configure
```

This will prompt you for:
- AWS Access Key ID
- AWS Secret Access Key
- Default region
- Output format

## Next Steps

Now that you have cfnfind installed, continue to the [Getting Started](./getting-started.md) tutorial to learn how to use it.
