# Terragrunt Shared Infrastructure Units 🧩

## Overview

This directory contains shared infrastructure unit configurations that provide a modular, flexible, and standardized approach to defining infrastructure components across different environments and modules.

## Architecture Principles 🏗️

### Core Design Concepts

1. **Modularity**: Discrete, self-contained infrastructure components
2. **Flexibility**: Environment-specific customizations with consistent base configuration
3. **Traceability**: Intelligent metadata and tagging management
4. **Dependency Orchestration**: Efficient cross-module dependency management

## Configuration Strategy 📦

### 1. Module Source Management

Implements a sophisticated, flexible module sourcing mechanism:

```hcl
locals {
  # Centralized module source configuration
  git_base_url           = "git::git@github.com:"
  tf_module_repository   = "your-org/terraform-modules"
  tf_module_version_default = "v0.1.0"
  tf_module_path_default = "modules/infrastructure-component"

  # Dynamic module source generation
  tf_module_source = format(
    "%s%s//%s",
    local.git_base_url,
    local.tf_module_repository,
    local.tf_module_path_default
  )
}
```

#### Source Management Benefits

- Centralized version control
- Consistent module referencing
- Flexible versioning
- Simplified update process

### 2. Intelligent Tagging System

Multi-layered tagging strategy for comprehensive resource management:

```hcl
locals {
  # Hierarchical Tagging Approach
  global_tags = {
    ManagedBy     = "Terragrunt"
    Architecture  = "Reference"
  }

  env_tags = {
    Environment = var.environment
    Region      = var.region
  }

  unit_tags = {
    Unit = "infrastructure-component"
    Type = "generator"
  }

  # Merged tags with clear precedence
  all_tags = merge(
    local.global_tags,
    local.env_tags,
    local.unit_tags
  )
}
```

#### Tagging Advantages

- Consistent resource identification
- Flexible tag management
- Enhanced tracking capabilities
- Cost allocation support
- Clear ownership definition

### 3. Dynamic Configuration Loading

Leverages Terragrunt's advanced configuration capabilities:

```hcl
locals {
  # Hierarchical configuration resolution
  global_config = read_terragrunt_config(find_in_parent_folders("global.hcl", ""))
  env_config    = read_terragrunt_config(find_in_parent_folders("env.hcl", ""))
  stack_config  = read_terragrunt_config(find_in_parent_folders("stack.hcl", ""))
}
```

#### Configuration Management Benefits

- Hierarchical configuration control
- Environment-specific customization
- Single source of truth
- Flexible override mechanisms

### 4. Dependency Management

Robust dependency resolution and mocking:

```hcl
dependency "prerequisite_module" {
  config_path = "../dependent-module"

  mock_outputs = {
    # Provides predictable outputs for planning
    module_output = "mock-value"
  }
}
```

#### Dependency Handling Advantages

- Explicit dependency declaration
- Validation-friendly mock outputs
- Controlled module interactions
- Simplified testing workflows

## Integration Mechanism 🔗

Child Terragrunt configurations include shared units using a standardized approach:

```hcl
include "shared" {
  path           = "${get_terragrunt_dir()}/../../../_shared/_units/component.hcl"
  expose         = true
  merge_strategy = "deep"
}
```

## Best Practices 🌟

### 1. Module Management

- Use semantic versioning
- Maintain stable module interfaces
- Document breaking changes
- Implement comprehensive input validation

### 2. Configuration Design

- Maintain focused, modular units
- Follow consistent naming conventions
- Implement clear tagging strategies
- Document configuration purpose and usage

### 3. Dependency Handling

- Define explicit, clear dependencies
- Provide meaningful mock outputs
- Document inter-module relationships
- Implement defensive configuration checks

## Troubleshooting 🛠️

### Common Resolution Strategies

1. **Module Source Verification**

   - Validate git repository access
   - Confirm module path accuracy
   - Check version constraint compatibility

2. **Configuration Loading Diagnostics**

   - Verify file path correctness
   - Validate environment variable configurations
   - Check configuration syntax and structure

3. **Dependency Resolution**
   - Validate dependency paths
   - Verify mock output completeness
   - Ensure output references are correct
