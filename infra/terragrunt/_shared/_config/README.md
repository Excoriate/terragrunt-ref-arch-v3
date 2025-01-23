# Terragrunt Shared Configuration Management 🛠️

## Overview

This directory contains centralized configuration files that provide a robust, flexible framework for managing infrastructure deployments across different environments and projects.

## Architecture Principles 🏗️

### Core Design Concepts

1. **Centralization**: Single source of truth for infrastructure configurations
2. **Flexibility**: Environment-aware, dynamically resolvable settings
3. **Traceability**: Comprehensive metadata and resource identification
4. **Consistency**: Uniform configuration management across infrastructure components

## Configuration Components 📦

### 1. Resource Tagging Strategy (`tags.hcl`)

Implements a sophisticated, multi-layered tagging mechanism:

```hcl
locals {
  # Hierarchical Tagging Approach
  global_tags = {
    ManagedBy     = "Terragrunt"
    OrchestratedBy = "Infrastructure-as-Code"
    Architecture  = "Reference"
  }

  environment_tags = {
    Environment = var.environment
    Region      = var.region
  }

  resource_tags = {
    Type    = "infrastructure-component"
    Project = local.project_name
    Version = local.project_version
  }

  # Merged tags with intelligent precedence
  all_tags = merge(
    local.global_tags,
    local.environment_tags,
    local.resource_tags
  )
}
```

#### Tagging Benefits

- **Resource Identification**: Precise, hierarchical resource tracking
- **Cost Allocation**: Granular resource categorization
- **Compliance Management**: Standardized metadata enforcement
- **Automated Governance**: Consistent resource labeling

### 2. Project Metadata (`app.hcl`)

Centralizes project-level configuration with dynamic resolution:

```hcl
locals {
  # Dynamic project metadata configuration
  project_name    = get_env("TG_STACK_APP_PRODUCT_NAME", "default-project")
  project_version = get_env("TG_STACK_APP_PRODUCT_VERSION", "0.0.0")
  environment     = get_env("TG_ENVIRONMENT", "development")

  # Computed project attributes
  project_identifier = lower(replace(local.project_name, "/[^a-zA-Z0-9]/", "-"))
}
```

#### Metadata Management Advantages

- Environment-driven configuration
- Flexible default value handling
- Normalized project identification
- Consistent metadata across infrastructure

### 3. Remote State Management (`remote_state.hcl`)

Implements a robust, secure remote state configuration strategy:

```hcl
locals {
  # Intelligent remote state configuration
  remote_state_bucket = get_env("TG_STACK_REMOTE_STATE_BUCKET_NAME", "")
  lock_table_name     = get_env("TG_STACK_REMOTE_STATE_LOCK_TABLE", "")
  state_region        = get_env("TG_STACK_REMOTE_STATE_REGION", "us-east-1")

  # State file naming strategy
  state_object_basename = get_env("TG_STACK_REMOTE_STATE_OBJECT_BASENAME", "terraform.tfstate")
  backend_filename      = get_env("TG_STACK_REMOTE_STATE_BACKEND_TF_FILENAME", "backend.tf")

  # Computed state configuration
  state_key = format(
    "%s/%s/%s/terraform.tfstate",
    local.project_identifier,
    local.environment,
    basename(get_terragrunt_dir())
  )
}
```

#### Remote State Management Benefits

- Dynamic, environment-driven configuration
- Secure, consistent state file naming
- Flexible region and bucket management
- Predictable state key generation

## Environment Variable Management 🌐

### Recommended Configuration

```bash
# Project Metadata
export TG_STACK_APP_PRODUCT_NAME="infrastructure-reference-arch"
export TG_STACK_APP_PRODUCT_VERSION="1.0.0"
export TG_ENVIRONMENT="production"

# Remote State Configuration
export TG_STACK_REMOTE_STATE_BUCKET_NAME="org-terraform-state"
export TG_STACK_REMOTE_STATE_LOCK_TABLE="terraform-state-locks"
export TG_STACK_REMOTE_STATE_REGION="us-east-1"
```

## Configuration Resolution Mechanism 🔄

### Dynamic Loading Strategy

1. Terragrunt recursively loads configurations from `_shared/_config`
2. Local variables computed using `read_terragrunt_config()`
3. Configurations merged with intelligent precedence
4. Environment variables override default values

### Inheritance and Override Patterns

- Shared configurations serve as base templates
- Module-specific configurations can extend or override shared settings
- Use `read_terragrunt_config()` for flexible configuration loading
- Implement merge strategies to control configuration inheritance

## Best Practices 🌟

1. **Naming Conventions**

   - Use lowercase, hyphen-separated project names
   - Follow semantic versioning
   - Maintain consistent naming across environments

2. **Security Considerations**

   - Never commit sensitive information
   - Use environment variables for dynamic configuration
   - Implement least-privilege access for state management

3. **Configuration Management**
   - Keep configurations DRY (Don't Repeat Yourself)
   - Document configuration purpose and usage
   - Validate configurations across different environments

## Troubleshooting 🛠️

### Common Resolution Strategies

1. **Configuration Validation**

   - Verify environment variable settings
   - Check configuration file syntax
   - Validate merge strategy implementations

2. **State Management Diagnostics**

   - Confirm S3 bucket and DynamoDB table existence
   - Validate IAM permissions
   - Check region and endpoint configurations

3. **Tagging Consistency**
   - Audit resource tags across infrastructure
   - Verify tag inheritance mechanisms
   - Implement automated tag validation
