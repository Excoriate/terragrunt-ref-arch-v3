# Environment Variables Management

## Overview

The Terragrunt Reference Architecture implements a sophisticated, hierarchical environment variable management system powered by [direnv](https://direnv.net/).

## Key Features

- **Hierarchical Inheritance**: Variables cascade from parent to child directories
- **Layer-based Organization**: Variables grouped by logical layers
- **Secure Variable Handling**: Validation and export of environment variables
- **Flexible Customization**: Easy-to-modify configuration points

## Configuration Hierarchy

```
/
├── .envrc                      # Root-level global configuration
└── infra/terragrunt/
    ├── .envrc                  # Terragrunt-specific variables
    ├── global/.envrc           # Global environment variables
    ├── dev/.envrc              # Development environment variables
    ├── staging/.envrc          # Staging environment variables
    └── prod/.envrc             # Production environment variables
```

## Configuration Principles

### Inheritance Mechanism

1. **Root Configuration**: Sets global defaults
2. **Terragrunt Layer**: Defines project-wide Terragrunt settings
3. **Environment Layers**: Provide environment-specific configurations
   - Each layer can override or extend parent configurations

### Core Utility Functions

- `_safe_export`: Securely export environment variables
- `_display_exported_vars`: Display current environment configuration
- `_log`: Standardized logging mechanism

## Example Configuration

### Root .envrc (Typical Configuration)

```bash
#!/usr/bin/env bash
# Exit immediately if any command fails
set -e

# Set project root
PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)}"
export PROJECT_ROOT

# Source utility functions
source "${PROJECT_ROOT}/scripts/envrc-utils.sh"

# Core initialization
_core_init

# Default configurations
_safe_export DEFAULT_REGION "us-east-1"
_safe_export TF_INPUT "0"
_safe_export LOG_LEVEL "info"

# Application Metadata
_safe_export TG_STACK_APP_PRODUCT_NAME "your-app-name"
```

### Environment-Specific .envrc (e.g., staging/.envrc)

```bash
#!/usr/bin/env bash
# Inherit from parent configuration
source_up || {
    echo >&2 "Warning: Could not source parent .envrc"
}

# Source utility functions
source "${PROJECT_ROOT}/scripts/envrc-utils.sh"

# Environment-Specific Variables
# Uncomment and modify as needed
# _safe_export AWS_PROFILE "staging"
# _safe_export TG_INSTANCE_TYPE "t3.medium"
# _safe_export TG_MAX_CAPACITY "4"

# Display exported variables
_display_exported_vars "STAGING_"
```

## Variable Customization

### Adding New Variables

1. **Global Variables**:
   ```bash
   _safe_export GLOBAL_SETTING "value"
   ```

2. **Environment-Specific Variables**:
   ```bash
   # In dev/.envrc
   _safe_export DEV_DEBUG_MODE "true"
   
   # In prod/.envrc
   _safe_export PROD_RESOURCE_SCALING "high"
   ```

## Best Practices

- Use `_safe_export` for all variable exports
- Leverage `source_up` for configuration inheritance
- Use environment-specific prefixes (DEV_, STAGING_, PROD_)
- Keep sensitive information out of version control

## Troubleshooting

### Common Issues

1. **Variables Not Loading**
   - Ensure `direnv` is installed
   - Run `direnv allow` in the directory
   - Check for syntax errors in `.envrc` files

2. **Inheritance Problems**
   - Verify `source_up` is present
   - Check for syntax errors preventing sourcing

### Debugging Commands

```bash
# Show all environment variables
env

# Show variables with specific prefix
env | grep TG_

# Direnv debug mode
DIRENV_LOG_FORMAT="" direnv allow
```

## Recommended Tools

- [direnv](https://direnv.net/): Environment variable management
- [sops](https://github.com/mozilla/sops): Secrets management
