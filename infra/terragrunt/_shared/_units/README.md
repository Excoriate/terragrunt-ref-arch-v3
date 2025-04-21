# Terragrunt Shared Infrastructure Units 🧩

## Overview

This directory contains shared infrastructure unit configurations that provide a modular, flexible, and standardized approach to defining infrastructure components across different environments and modules.

```txt
infra/terragrunt/
├── _shared/
│   ├── _config/
│   ├── _units/
│   │   └── <unit>.hcl # unit referenced in the consuming unit's terragrunt.hcl (path pattern: infra/terragrunt/<env>/<stack>/<unit>/terragrunt.hcl)
│   └── README.md
```


## Architecture Principles 🏗️

### Core Design Concepts

1. **Modularity**: Discrete, self-contained infrastructure components
2. **Flexibility**: Environment-specific customizations with consistent base configuration
3. **Traceability**: Intelligent metadata and tagging management
4. **Dependency Orchestration**: Efficient cross-module dependency management

## Configuration Strategy 📦

### 1. Module Source Management

Automatically reads the environment configuration:

```hcl
# file located in infra/terragrunt/<env>/env.hcl
env_cfg = read_terragrunt_config(find_in_parent_folders("env.hcl"))
```

Automatically reads the stack configuration:

```hcl
# file located in infra/terragrunt/<env>/<stack>/stack.hcl
stack_cfg = read_terragrunt_config(find_in_parent_folders("stack.hcl"))
```

Automatically reads the global configuration from the root `config.hcl` file:

```hcl
# file located in infra/terragrunt/config.hcl
cfg = read_terragrunt_config("${find_in_parent_folders("config.hcl")}")
```

### 2. Intelligent Tagging System

The tags are divided into three distinct layers:

1. **Global Tags**: Tags that are applied to all resources.
2. **Environment Tags**: Tags that are applied to all resources in the environment.
3. **Unit Tags**: Tags that are applied to all resources in the unit.

Each unit is also able to decorate with its own tags, like so:

```hcl
unit_tags = {
  Unit = "infrastructure-component"
  Type = "generator"
}
```

At the end, these tags are merged together, with a clear precedence strategy:

```hcl
  # 🔗 TAG SOURCE AGGREGATION
  # Collect tags from different configuration levels
  env_tags    = local.env_cfg.locals.tags
  global_tags = local.cfg.locals.tags
  stack_tags  = local.stack_cfg.locals.tags

  # 🧩 FINAL TAG COMPOSITION
  # Merge tags with a clear precedence strategy
  # Precedence: Unit Tags > Environment Tags > Global Tags
  all_tags = merge(
    local.env_tags,
    local.unit_tags,
    local.global_tags,
    local.stack_tags
  )
```

This way, each unit can be tagged with its own tags, while still benefiting from the global tags and environment tags.

### 3. Git Module Source Management

The module source is managed through a centralized configuration file in the `_config` directory (see [README.md](../_config/README.md) for more details).

| Attribute | Description | Default Value | Environment Variable Override |
|-----------|-------------|---------------|-------------------------------|
| `git_base_url` | Base URL for GitHub repository access | `local.cfg.locals.cfg_git.git_base_urls.github` | N/A |
| `git_base_url_local` | Base URL for local repository access | `local.cfg.locals.cfg_git.git_base_urls.local` | N/A |
| `tf_module_repository` | GitHub repository path for the Terraform module | `"Excoriate/terraform-aws-codeartifact"` | N/A |
| `tf_module_path_default` | Default path within the repository where the module is located | `"modules/domain-permissions"` | N/A |
| `tf_module_version_default` | Default version tag of the module | `"v1.1.2"` | `TG_STACK_TF_MODULE_<MODULE_NAME>_VERSION_OVERRIDE` |

This configuration allows for quick iterations, and local overrides for troubleshooting and debugging.

>NOTE: On each run, the module's version is echoed to the console for visibility.

