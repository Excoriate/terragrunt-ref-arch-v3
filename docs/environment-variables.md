# Environment Variables Management

## Overview

The Terragrunt Reference Architecture implements a simplified yet powerful environment variable management system powered by [direnv](https://direnv.net/). This architecture now uses a single root `.envrc` file with category-based organization for improved clarity and easier customization.

## Key Features

- **Category-Based Organization**: Variables grouped by functional categories
- **Secure Variable Handling**: Validation and export of environment variables
- **Flexible Customization**: Dedicated section for custom user variables
- **Visual Clarity**: Emoji-tagged sections for improved readability

## Configuration Structure

```
/
└── .envrc                      # Comprehensive root-level configuration
```

## Configuration Principles

### Organizational Structure

The root `.envrc` file is organized into the following categories:

1. **Project Metadata**: Core project information and authorship
2. **Cloud Provider Settings**: Region configuration and provider-specific settings
3. **Terraform & Terragrunt Configuration**: Tool versions and behavior settings
4. **Logging & Diagnostics**: Output verbosity and log storage
5. **Remote State Configuration**: Backend storage for Terraform state
6. **Custom Use-Case Variables**: User-defined environment variables

### Core Utility Functions

- `_safe_export`: Securely export environment variables
- `_layer_export`: Export variables with additional layer-specific logging
- `_display_exported_vars`: Display current environment configuration
- `_log`: Standardized logging mechanism
- `_layer_env_info`: Display organized layer information with descriptions

## Root .envrc (Comprehensive Configuration)

```bash
#!/usr/bin/env bash
# Terragrunt Reference Architecture - Environment Configuration
# Simplified and modular environment setup

# Exit immediately if any command fails
set -e

# Ensure PROJECT_ROOT is set reliably
PROJECT_ROOT="${PROJECT_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)}"
export PROJECT_ROOT

# Source utility functions
source "${PROJECT_ROOT}/scripts/envrc-utils.sh"

# Core initialization
_core_init

# =====================================================================
# 🔧 CUSTOMIZATION SECTION
# =====================================================================
# Configuration variables are grouped by functional categories for easier
# customization and maintenance.

# ---------------------------------------------------------------------
# 1️⃣ PROJECT METADATA
# ---------------------------------------------------------------------
# Define core project information and authorship
# ---------------------------------------------------------------------
TG_STACK_APP_AUTHOR="${TG_STACK_APP_AUTHOR:-Your Name}"
_safe_export TG_STACK_APP_AUTHOR "$TG_STACK_APP_AUTHOR"

TG_STACK_APP_PRODUCT_NAME="${TG_STACK_APP_PRODUCT_NAME:-your-app-name}"
_safe_export TG_STACK_APP_PRODUCT_NAME "$TG_STACK_APP_PRODUCT_NAME"

# ---------------------------------------------------------------------
# 2️⃣ CLOUD PROVIDER & REGION SETTINGS
# ---------------------------------------------------------------------
# Configure cloud provider-specific settings
# ---------------------------------------------------------------------
DEFAULT_REGION="${DEFAULT_REGION:-us-east-1}"
_safe_export DEFAULT_REGION "$DEFAULT_REGION"

# ---------------------------------------------------------------------
# 3️⃣ TERRAFORM & TERRAGRUNT CONFIGURATION
# ---------------------------------------------------------------------
# Control Terraform behavior and version requirements
# ---------------------------------------------------------------------
# Core Terraform Settings
TF_INPUT="${TF_INPUT:-0}"
_safe_export TF_INPUT "$TF_INPUT"

TG_STACK_TF_VERSION="${TG_STACK_TF_VERSION:-1.9.0}"
_safe_export TG_STACK_TF_VERSION "$TG_STACK_TF_VERSION"

# Terragrunt Performance Settings
TERRAGRUNT_DOWNLOAD_DIR="${TERRAGRUNT_DOWNLOAD_DIR:-${HOME}/.terragrunt-cache/$(basename "${PROJECT_ROOT}")}"
_safe_export TERRAGRUNT_DOWNLOAD_DIR "$TERRAGRUNT_DOWNLOAD_DIR"

TERRAGRUNT_CACHE_MAX_AGE="${TERRAGRUNT_CACHE_MAX_AGE:-168h}"
_safe_export TERRAGRUNT_CACHE_MAX_AGE "$TERRAGRUNT_CACHE_MAX_AGE"

# Terragrunt Behavior Settings
TERRAGRUNT_LOG_LEVEL="${TERRAGRUNT_LOG_LEVEL:-info}"
_safe_export TERRAGRUNT_LOG_LEVEL "$TERRAGRUNT_LOG_LEVEL"

TERRAGRUNT_AUTO_INIT="${TERRAGRUNT_AUTO_INIT:-true}"
_safe_export TERRAGRUNT_AUTO_INIT "$TERRAGRUNT_AUTO_INIT"

# ---------------------------------------------------------------------
# 6️⃣ CUSTOM USE-CASE VARIABLES
# ---------------------------------------------------------------------
# Add your custom environment variables below
# Examples:
# _safe_export TG_CUSTOM_VAR_1 "value1"
# _safe_export TG_CUSTOM_VAR_2 "value2"
# ---------------------------------------------------------------------
```

## Variable Customization

### Adding New Custom Variables

To add your own custom variables, locate the "CUSTOM USE-CASE VARIABLES" section in the root `.envrc` file and add your variables there:

```bash
# ---------------------------------------------------------------------
# 6️⃣ CUSTOM USE-CASE VARIABLES
# ---------------------------------------------------------------------
# Add your custom environment variables below

# Development-specific settings
_safe_export TG_DEV_DEBUG_MODE "true"
_safe_export TG_DEV_API_ENDPOINT "https://dev-api.example.com"

# Production-specific settings
_safe_export TG_PROD_RESOURCE_SCALING "high"
_safe_export TG_PROD_REPLICA_COUNT "5"
```

## Best Practices

- Use descriptive variable names with appropriate prefixes
- Group related variables together in the custom section
- Use `_safe_export` for all variable exports to ensure proper validation
- Add comments to document the purpose of custom variables
- Keep sensitive information out of version control

## Troubleshooting

### Common Issues

1. **Variables Not Loading**
   - Ensure `direnv` is installed: `which direnv`
   - Run `direnv allow` in the directory
   - Check for syntax errors in the `.envrc` file

2. **Validation Failures**
   - Verify required variables are properly defined
   - Check the output of `_validate_layer_config` in the logs

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
