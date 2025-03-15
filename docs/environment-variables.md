# Environment Variables Guide

## Overview

This project uses [direnv](https://direnv.net/) with `.envrc` files to manage environment variables. This approach provides a secure, consistent way to handle configuration across different environments and developers.

## Benefits of Using `.envrc` with direnv

1. **Automatic Loading**: Environment variables are loaded automatically when entering the project directory
2. **Secure Management**: Variables are never committed to version control
3. **Project-Specific Configuration**: Each project can have its own isolated environment
4. **Shell Integration**: Works with bash, zsh, and other common shells

## Setup Instructions

### Prerequisites

1. Install direnv:
   ```bash
   # macOS with Homebrew
   brew install direnv
   
   # Linux
   sudo apt install direnv  # Debian/Ubuntu
   sudo yum install direnv  # CentOS/RHEL
   ```

2. Add direnv hook to your shell:
   ```bash
   # For bash
   echo 'eval "$(direnv hook bash)"' >> ~/.bashrc
   source ~/.bashrc
   
   # For zsh
   echo 'eval "$(direnv hook zsh)"' >> ~/.zshrc
   source ~/.zshrc
   ```

### Project Configuration

1. Create or edit the `.envrc` file in the project root:
   ```bash
   # Example: Setting up basic environment variables
   export TG_STACK_REGION="us-east-1"
   export TG_STACK_APP_PRODUCT_NAME="my-infrastructure"
   ```

2. Allow direnv to load the environment:
   ```bash
   direnv allow
   ```

3. Verify the environment is loaded:
   ```bash
   echo $TG_STACK_REGION  # Should output "us-east-1"
   ```

## Required Environment Variables

The following environment variables are used throughout the project:

### Application Observability
- `LOG_LEVEL`: Minimum log level for application logging
- `LOG_DIR`: Directory where application logs will be stored

### Terraform Remote State Management
- `TG_STACK_REMOTE_STATE_BUCKET_NAME`: S3 bucket for storing Terraform state
- `TG_STACK_REMOTE_STATE_LOCK_TABLE`: DynamoDB table for state locking
- `TG_STACK_REMOTE_STATE_REGION`: AWS region for remote state storage
- `TG_STACK_REMOTE_STATE_OBJECT_BASENAME`: Base filename for state objects

### Terragrunt Configuration
- `TG_STACK_FLAG_ENABLE_PROVIDERS_OVERRIDE`: Enable provider configuration overrides
- `TG_STACK_FLAG_ENABLE_VERSIONS_OVERRIDE`: Enable version configuration overrides
- `TG_STACK_REGION`: Default AWS deployment region
- `TG_STACK_TF_VERSION`: Enforced Terraform version
- `TG_STACK_FLAG_ENABLE_TERRAFORM_VERSION_FILE_OVERRIDE`: Enable Terraform version file override

### Application Metadata
- `TG_STACK_APP_AUTHOR`: Author of the application
- `TG_STACK_APP_PRODUCT_NAME`: Name of the application

## Example `.envrc` Configuration

Here's a complete example of an `.envrc` file with all required variables:

```bash
#!/usr/bin/env bash
# Project environment configuration

# Enable Nix flake support for direnv (if using Nix)
use flake

# Global defaults and security settings
export DEFAULT_REGION="us-east-1"
export TF_INPUT="0"  # Disable interactive Terraform input
export LANG="en_US.UTF-8"
export LC_ALL="en_US.UTF-8"

# Application Observability Configuration
export LOG_LEVEL="info"
export LOG_DIR="/var/log/myapp"

# Terraform Remote State Management
export TG_STACK_REMOTE_STATE_BUCKET_NAME="terraform-state-myproject"
export TG_STACK_REMOTE_STATE_LOCK_TABLE="terraform-state-lock-myproject"
export TG_STACK_REMOTE_STATE_REGION="us-east-1"
export TG_STACK_REMOTE_STATE_OBJECT_BASENAME="terraform.tfstate.json"

# Terragrunt Configuration Variables
export TG_STACK_FLAG_ENABLE_PROVIDERS_OVERRIDE="true"
export TG_STACK_FLAG_ENABLE_VERSIONS_OVERRIDE="true"
export TG_STACK_REGION="us-east-1"
export TG_STACK_TF_VERSION="1.9.0"
export TG_STACK_FLAG_ENABLE_TERRAFORM_VERSION_FILE_OVERRIDE="false"

# Application Metadata
export TG_STACK_APP_AUTHOR="YourName"
export TG_STACK_APP_PRODUCT_NAME="my-infrastructure"
```

## Customizing Environment Variables

To customize environment variables for your local setup:

1. Copy the example above to your project's `.envrc` file
2. Modify the values to match your environment
3. Run `direnv allow` to apply the changes

## Environment-Specific Configurations

For different environments (development, staging, production):

1. Create environment-specific `.envrc` files in subdirectories
2. Use `source_up` in subdirectory `.envrc` files to inherit parent variables
3. Override specific variables as needed

Example for a staging environment:

```bash
# staging/.envrc
source_up
export TG_STACK_REGION="us-west-2"
export TG_STACK_APP_PRODUCT_NAME="my-infrastructure-staging"
```

## Troubleshooting

If you encounter issues with environment variables not being loaded:

1. Ensure direnv is properly installed and hooked into your shell
2. Run `direnv allow` in the project directory
3. Check for any errors in the `.envrc` file
4. Verify that your shell is properly configured to use direnv

For more detailed information about direnv, visit [direnv.net](https://direnv.net/) 
