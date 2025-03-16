# Environment Variables Management

## Overview

The Terragrunt Reference Architecture employs a sophisticated, hierarchical approach to environment variable management powered by [direnv](https://direnv.net/). This system provides:

- **Hierarchical Inheritance**: Variables cascade from parent to child directories
- **Layer-based Organization**: Variables are grouped by logical layers
- **Secure Variable Handling**: Validation and export of environment variables
- **Flexible Customization**: Easy-to-modify configuration points
- **Comprehensive Utility Functions**: Shared shell functions for variable management

## Core Principles

### Inheritance Mechanism

The environment variable system is built on a cascading inheritance model:

1. **Root Level**: Global project-wide defaults
2. **Terragrunt Layer**: Terragrunt-specific configurations
3. **Environment Layers**: Environment-specific settings (dev, staging, prod)
4. **Optional Stack Layers**: Granular, stack-specific configurations

### Key Utility Functions

The system leverages several core utility functions in `scripts/envrc-utils.sh`:

- `_safe_export`: Securely export environment variables
- `_display_exported_vars`: Show current environment variable configuration
- `_detect_project_root`: Identify project root directory
- `_log`: Standardized logging mechanism

## Directory Structure

```
/
├── .envrc                      # Root-level global variables
├── scripts/
│   └── envrc-utils.sh          # Shared utility functions
├── infra/
│   └── terragrunt/
│       ├── .envrc              # Terragrunt-specific variables
│       ├── global/
│       │   └── .envrc          # Global environment variables
│       ├── dev/
│       │   └── .envrc          # Development environment variables
│       ├── staging/
│       │   └── .envrc          # Staging environment variables
│       └── prod/
│           └── .envrc          # Production environment variables
```

## `.envrc` File Template

Each environment's `.envrc` follows a consistent structure:

```bash
#!/usr/bin/env bash
# Environment-Specific Configuration

# Inherit from parent configuration
source_up || true

# Source utility functions
source "${PROJECT_ROOT}/scripts/envrc-utils.sh"

# Environment Variables
# Uncomment and modify as needed:
# _safe_export ENVIRONMENT_TYPE "specific-type"
# _safe_export RESOURCE_LIMITS "configuration"

# Optional: Display exported variables
_display_exported_vars "ENV_PREFIX_"
```

## Customization Strategies

### Adding New Variables

1. **Global Variables** (root `.envrc`):
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

### Variable Inheritance and Overriding

- Variables defined in child directories override parent configurations
- Use `source_up` to inherit parent variables
- Customize by uncommenting and modifying variables

## Utility Functions Reference

### `_safe_export`
Securely export environment variables with validation:
```bash
_safe_export VARIABLE_NAME "value"
```

### `_display_exported_vars`
Show exported variables with optional prefix filtering:
```bash
_display_exported_vars           # Show all variables
_display_exported_vars "TG_"     # Show variables starting with TG_
```
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

## Advanced Configuration

### Conditional Configurations

```bash
# Example of environment-specific configuration
if [[ "$TG_ENV" == "dev" ]]; then
  _safe_export DEV_SPECIFIC_SETTING "value"
fi
```

### Custom Validation

```bash
# Custom validation function
_validate_environment() {
  local env_type="$1"
  if [[ ! "$env_type" =~ ^(dev|staging|prod)$ ]]; then
    _log "ERROR" "Invalid environment type: $env_type"
    return 1
  fi
}
```

## Security Considerations

- Variables containing "SECRET", "PASSWORD", or "KEY" are automatically masked
- All variables are validated before export
- Caller information is tracked for each variable set

## Integration with Other Tools

- Compatible with Terraform
- Works seamlessly with Terragrunt
- Supports Nix flake configurations

## Recommended Tools

- [direnv](https://direnv.net/): Environment variable management
- [sops](https://github.com/mozilla/sops): Secrets management
- [pre-commit](https://pre-commit.com/): Git hooks for validation
