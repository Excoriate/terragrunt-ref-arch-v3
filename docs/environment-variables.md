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

1.  **Project Metadata**: Core project information and authorship
2.  **Cloud Provider Settings**: Region configuration and provider-specific settings
3.  **Terraform & Terragrunt Configuration**: Tool versions and behavior settings
4.  **Logging & Diagnostics**: Output verbosity and log storage
5.  **Remote State Configuration**: Backend storage for Terraform state
6.  **Custom Use-Case Variables**: User-defined environment variables

### Core Utility Functions

- `_safe_export`: Securely export environment variables
- `_layer_export`: Export variables with additional layer-specific logging
- `_display_exported_vars`: Display current environment configuration
- `_log`: Standardized logging mechanism
- `_layer_env_info`: Display organized layer information with descriptions

## Root .envrc (Comprehensive Configuration)

The following is the content of the root [`.envrc`](../../.envrc) file, which serves as the central point for environment variable definition:

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
# shellcheck source=./scripts/envrc-utils.sh
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
# The following uses the new recommended environment variables
# instead of the deprecated ones

# Replace deprecated TF_INPUT with TG_NON_INTERACTIVE
TG_NON_INTERACTIVE="${TG_NON_INTERACTIVE:-true}"
_safe_export TG_NON_INTERACTIVE "${TG_NON_INTERACTIVE}"

TG_STACK_TF_VERSION="${TG_STACK_TF_VERSION:-1.9.0}"
_safe_export TG_STACK_TF_VERSION "$TG_STACK_TF_VERSION"

# Terragrunt Performance Settings
TERRAGRUNT_DOWNLOAD_DIR="${TERRAGRUNT_DOWNLOAD_DIR:-${HOME}/.terragrunt-cache/$(basename "${PROJECT_ROOT}")}"
_safe_export TERRAGRUNT_DOWNLOAD_DIR "$TERRAGRUNT_DOWNLOAD_DIR"

TERRAGRUNT_CACHE_MAX_AGE="${TERRAGRUNT_CACHE_MAX_AGE:-168h}"
_safe_export TERRAGRUNT_CACHE_MAX_AGE "$TERRAGRUNT_CACHE_MAX_AGE"

# Terragrunt Behavior Settings
# Replace deprecated TERRAGRUNT_LOG_LEVEL with TG_LOG_LEVEL
TG_LOG_LEVEL="${TG_LOG_LEVEL:-info}"
_safe_export TG_LOG_LEVEL "${TG_LOG_LEVEL}"

TERRAGRUNT_DISABLE_CONSOLE_OUTPUT="${TERRAGRUNT_DISABLE_CONSOLE_OUTPUT:-false}"
_safe_export TERRAGRUNT_DISABLE_CONSOLE_OUTPUT "${TERRAGRUNT_DISABLE_CONSOLE_OUTPUT}"

# Replace deprecated TERRAGRUNT_AUTO_INIT with TG_NO_AUTO_INIT
# Note: The logic is inverted, so we use false when auto-init is desired
TG_NO_AUTO_INIT="${TG_NO_AUTO_INIT:-false}"
_safe_export TG_NO_AUTO_INIT "${TG_NO_AUTO_INIT}"

# Replace deprecated TERRAGRUNT_AUTO_RETRY with TG_NO_AUTO_RETRY
# Note: The logic is inverted, so we use false when auto-retry is desired
TG_NO_AUTO_RETRY="${TG_NO_AUTO_RETRY:-false}"
_safe_export TG_NO_AUTO_RETRY "${TG_NO_AUTO_RETRY}"

# ---------------------------------------------------------------------
# 4️⃣ LOGGING & DIAGNOSTICS
# ---------------------------------------------------------------------
# Configure output verbosity and log storage
# ---------------------------------------------------------------------
LOG_LEVEL="${LOG_LEVEL:-info}"
_safe_export LOG_LEVEL "${LOG_LEVEL}"

LOG_DIR_PATH="${HOME}/.logs/${TG_STACK_APP_PRODUCT_NAME}"
# Ensure log directory exists
mkdir -p "${LOG_DIR_PATH}" 2>/dev/null || true
_safe_export LOG_DIR "${LOG_DIR_PATH}"

# ---------------------------------------------------------------------
# 5️⃣ REMOTE STATE CONFIGURATION
# ---------------------------------------------------------------------
# Define backend storage for Terraform state
# ---------------------------------------------------------------------
# Placeholder values - MUST be replaced in actual configuration
_safe_export TG_STACK_REMOTE_STATE_BUCKET_NAME "terraform-state-makemyinfra"
_safe_export TG_STACK_REMOTE_STATE_LOCK_TABLE "terraform-state-makemyinfra"

# ---------------------------------------------------------------------
# 6️⃣ CUSTOM USE-CASE VARIABLES
# ---------------------------------------------------------------------
# Add your custom environment variables below
# Examples:
# _safe_export TG_CUSTOM_VAR_1 "value1"
# _safe_export TG_CUSTOM_VAR_2 "value2"
# ---------------------------------------------------------------------

# Tool Availability Check
_validate_layer_config "CORE" \
  "DEFAULT_REGION" \
  "TG_STACK_APP_PRODUCT_NAME"

# Final initialization log
_log "INFO" "Environment for ${TG_STACK_APP_PRODUCT_NAME} initialized successfully"

# Display layer information for core Terragrunt settings
_layer_env_info "TERRAGRUNT" \
  "TERRAGRUNT_DOWNLOAD_DIR:Terragrunt Download Directory" \
  "TG_LOG_LEVEL:Terragrunt Log Level" \
  "TG_NO_AUTO_INIT:Disable Auto Initialize" \
  "TG_NO_AUTO_RETRY:Disable Auto Retry"

# Display all exported environment variables
_display_exported_vars ""
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

### Comprehensive Environment Variable List

This section provides a comprehensive list of environment variables used throughout the Terragrunt Reference Architecture, including those defined in the root `.envrc` file and those utilized within HCL configurations.

| Category                    | Variable Name                                       | Description                                                                 | Default Value (`.envrc` or HCL)           | Primary Source/Usage                                  |
|-----------------------------|-----------------------------------------------------|-----------------------------------------------------------------------------|-------------------------------------------|-------------------------------------------------------|
| **Project & Metadata**      | `PROJECT_ROOT`                                      | Absolute path to the project root directory.                                | Dynamically set in `.envrc`               | `.envrc`                                              |
|                             | `TG_STACK_APP_AUTHOR`                               | Author of the configuration.                                                | `Your Name`                               | `.envrc`, `infra/terragrunt/_shared/_config/tags.hcl` |
|                             | `TG_STACK_APP_PRODUCT_NAME`                         | Project/application name for identification.                                | `your-app-name`                           | `.envrc`, `infra/terragrunt/_shared/_config/app.hcl`, `infra/terragrunt/_shared/_config/tags.hcl` |
| **Cloud Provider & Region** | `DEFAULT_REGION`                                    | Default cloud provider region.                                              | `us-east-1`                               | `.envrc`                                              |
|                             | `TG_STACK_DEPLOYMENT_REGION`                        | AWS region for deployments, overrides `DEFAULT_REGION` for Terragrunt.    | `us-east-1` (from HCL)                    | `infra/terragrunt/config.hcl`                         |
| **Terraform & Terragrunt**  |                                                     |                                                                             |                                           |                                                       |
| *General*                   | `TG_NON_INTERACTIVE`                                | Replaces `TF_INPUT`. If true, disables interactive prompts.                 | `true`                                    | `.envrc`                                              |
|                             | `TG_STACK_TF_VERSION`                               | Enforced Terraform version.                                                 | `1.9.0` (`.envrc`), `1.11.3` (HCL)        | `.envrc`, `infra/terragrunt/config.hcl`               |
| *Terragrunt Performance*    | `TERRAGRUNT_DOWNLOAD_DIR`                           | Directory where Terragrunt downloads remote sources.                        | `${HOME}/.terragrunt-cache/...`            | `.envrc` (Terragrunt Standard)                        |
|                             | `TERRAGRUNT_CACHE_MAX_AGE`                          | Max age for items in Terragrunt cache.                                      | `168h`                                    | `.envrc` (Terragrunt Standard)                        |
| *Terragrunt Behavior*       | `TG_LOG_LEVEL`                                      | Terragrunt logging verbosity. Replaces `TERRAGRUNT_LOG_LEVEL`.            | `info`                                    | `.envrc` (Terragrunt Standard)                        |
|                             | `TERRAGRUNT_DISABLE_CONSOLE_OUTPUT`                 | If true, disables Terragrunt console output.                                | `false`                                   | `.envrc` (Terragrunt Standard)                        |
|                             | `TG_NO_AUTO_INIT`                                   | Replaces `TERRAGRUNT_AUTO_INIT` (inverted). If true, disables auto-init.   | `false` (meaning auto-init is enabled)    | `.envrc` (Terragrunt Standard)                        |
|                             | `TG_NO_AUTO_RETRY`                                  | Replaces `TERRAGRUNT_AUTO_RETRY` (inverted). If true, disables auto-retry. | `false` (meaning auto-retry is enabled)   | `.envrc` (Terragrunt Standard)                        |
| *Terragrunt Flags (HCL)*    | `TG_STACK_FLAG_ENABLE_PROVIDERS_OVERRIDE`           | Controls dynamic provider file generation.                                  | `true` (from HCL)                         | `infra/terragrunt/config.hcl`                         |
|                             | `TG_STACK_FLAG_ENABLE_VERSIONS_OVERRIDE`            | Controls dynamic version file generation.                                   | `true` (from HCL)                         | `infra/terragrunt/config.hcl`                         |
|                             | `TG_STACK_FLAG_ENABLE_TERRAFORM_VERSION_FILE_OVERRIDE` | Controls `.terraform-version` file generation.                             | `false` (from HCL)                        | `infra/terragrunt/config.hcl`                         |
| **Logging & Diagnostics**   | `LOG_LEVEL`                                         | General logging level for scripts.                                          | `info`                                    | `.envrc`                                              |
|                             | `LOG_DIR`                                           | Directory for storing logs.                                                 | `${HOME}/.logs/...` (derived)             | `.envrc`                                              |
| **Remote State Config**     | `TG_STACK_REMOTE_STATE_BUCKET_NAME`                 | S3 bucket for Terraform remote state.                                       | `terraform-state-makemyinfra`             | `.envrc`, `infra/terragrunt/_shared/_config/remote_state.hcl` |
|                             | `TG_STACK_REMOTE_STATE_LOCK_TABLE`                  | DynamoDB table for state locking.                                           | `terraform-state-makemyinfra`             | `.envrc`, `infra/terragrunt/_shared/_config/remote_state.hcl` |
|                             | `TG_STACK_REMOTE_STATE_REGION`                      | AWS region for remote state storage.                                        | `us-east-1` (from HCL)                    | `infra/terragrunt/_shared/_config/remote_state.hcl`   |
|                             | `TG_STACK_REMOTE_STATE_OBJECT_BASENAME`             | Basename for the Terraform state object file.                               | `terraform.tfstate.json` (from HCL)       | `infra/terragrunt/_shared/_config/remote_state.hcl`   |
|                             | `TG_STACK_REMOTE_STATE_BACKEND_TF_FILENAME`         | Filename for the generated backend configuration.                           | `backend.tf` (from HCL)                   | `infra/terragrunt/_shared/_config/remote_state.hcl`   |
| **Module Version Overrides**| `TG_STACK_TF_MODULE_DNI_GENERATOR_VERSION_DEFAULT`  | Default version for the DNI generator module.                               | `v0.1.0` (from HCL)                       | `infra/terragrunt/_shared/_units/dni_generator.hcl`   |
|                             | `TG_STACK_TF_MODULE_NAME_GENERATOR_VERSION_DEFAULT` | Default version for the name generator module.                              | `v0.1.0` (from HCL)                       | `infra/terragrunt/_shared/_units/name_generator.hcl`  |
|                             | `TG_STACK_TF_MODULE_AGE_GENERATOR_VERSION_DEFAULT`  | Default version for the age generator module.                               | `v0.1.0` (from HCL)                       | `infra/terragrunt/_shared/_units/age_generator.hcl`   |
|                             | `TG_STACK_TF_MODULE_LASTNAME_GENERATOR_VERSION_DEFAULT` | Default version for the lastname generator module.                          | `v0.1.0` (from HCL)                       | `infra/terragrunt/_shared/_units/lastname_generator.hcl`|

*Note: Some Terragrunt standard variables like `TF_INPUT`, `TERRAGRUNT_LOG_LEVEL`, `TERRAGRUNT_AUTO_INIT`, `TERRAGRUNT_AUTO_RETRY` have been replaced by their `TG_` prefixed counterparts (e.g., `TG_NON_INTERACTIVE`, `TG_LOG_LEVEL`, `TG_NO_AUTO_INIT`, `TG_NO_AUTO_RETRY`) in the root `.envrc` for consistency and to align with newer Terragrunt practices where applicable. The `TG_NO_AUTO_INIT` and `TG_NO_AUTO_RETRY` variables have inverted logic compared to their predecessors.*

## Best Practices

- Use descriptive variable names with appropriate prefixes
- Group related variables together in the custom section
- Use `_safe_export` for all variable exports to ensure proper validation
- Add comments to document the purpose of custom variables
- Keep sensitive information out of version control

## Troubleshooting

### Common Issues

1.  **Variables Not Loading**
    - Ensure `direnv` is installed: `which direnv`
    - Run `direnv allow` in the directory
    - Check for syntax errors in the `.envrc` file

2.  **Validation Failures**
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
