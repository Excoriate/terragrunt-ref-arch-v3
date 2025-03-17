# Infrastructure Templates 🧩

## Overview

This directory contains template files used across the infrastructure deployment pipeline, providing standardized configurations and version management.

## Current Templates

### `.terraform-version.tpl`

A version constraint template for Terraform, ensuring consistent tooling across different environments and team members.

#### Usage Example

```hcl
# Terraform version constraint
terraform {
  required_version = file("${path_relative_to_include()}/../_templates/.terraform-version.tpl")
}
```

## Purpose

- **Standardization**: Maintain consistent configuration across projects
- **Version Control**: Centralize version management for critical tools
- **Flexibility**: Easy to update and propagate changes

This template can be used in the `generate` block of a Terragrunt configuration file to ensure that the correct version of Terraform is used. See [root.hcl](../root.hcl) for an example.

# Environment Configuration Templates

This directory contains templates for environment configuration files that can be used across the project.

## Layer-Agnostic Environment Configuration

The environment configuration system is designed to be:

1. **Reusable** - Common functions are extracted to shared utility scripts
2. **Flexible** - Not tied to specific environment naming conventions
3. **Hierarchical** - Each layer inherits from its parent
4. **Self-documenting** - Includes clear documentation and logging

## Available Templates

- `layer.envrc.template` - Template for layer-specific `.envrc` files

## How to Use

1. Copy the template to your new layer directory:
   ```bash
   cp _templates/layer.envrc.template your/new/layer/.envrc
   ```

2. Customize the layer name and variables:
   ```bash
   # Change this to match your layer
   LAYER_NAME="YOUR_LAYER_NAME"
   ```

3. Add your layer-specific variables using the `_layer_export` function:
   ```bash
   _layer_export YOUR_VAR "your-value" "$LAYER_NAME"
   ```

4. Update the layer information display:
   ```bash
   _layer_env_info "$LAYER_NAME" \
     "YOUR_VAR:Description of YOUR_VAR" \
     "ANOTHER_VAR:Description of ANOTHER_VAR"
   ```

## Inheritance

Each `.envrc` file automatically inherits from its parent directory's `.envrc` file through the `source_up` directive. This creates a hierarchical configuration where:

1. Root `.envrc` sets global defaults
2. Each layer can override or extend these defaults
3. Nested layers inherit from their parent layers

## Utility Functions

The shared utility script (`scripts/envrc-utils.sh`) provides several helpful functions:

- `_layer_export` - Export a variable with layer-specific logging
- `_layer_log` - Log a message with layer-specific prefix
- `_layer_env_info` - Display layer environment information
- `_validate_layer_config` - Validate required variables for a layer

## Example Layer Structure

```
project/
├── .envrc                  # Root environment variables
├── scripts/
│   └── envrc-utils.sh      # Shared utility functions
└── infra/
    └── terragrunt/
        ├── .envrc          # Terragrunt-specific variables
        ├── dev/
        │   └── .envrc      # Development layer variables
        ├── staging/
        │   └── .envrc      # Staging layer variables
        └── prod/
            └── .envrc      # Production layer variables
```

Each layer can be named according to your specific needs - the system is not tied to specific environment names.
