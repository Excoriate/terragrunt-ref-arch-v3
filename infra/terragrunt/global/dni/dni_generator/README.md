# Infrastructure Unit Configuration 🛠️

## Overview

This Terragrunt unit demonstrates a dynamic provider and version management system implemented in the infrastructure. It provides a flexible, modular approach to configuring infrastructure components.

## 📁 File Structure

```
dns-zone/
├── unit_cfg_providers.hcl  # Provider configurations (optional)
├── unit_cfg_versions.hcl   # Version constraints (optional)
├── terragrunt.hcl          # Unit-specific Terragrunt configuration
└── README.md              # This documentation
```

## 🔌 Provider Configuration Management

### Dynamic Configuration System

The unit implements a flexible provider configuration mechanism:

1. **Local Provider Configuration** (`unit_cfg_providers.hcl`):

   - Defines provider-specific settings
   - Credentials sourced from environment variables
   - Supports multiple provider configurations

2. **Version Management** (`unit_cfg_versions.hcl`):
   - Specifies provider version constraints
   - Ensures consistent provider versions across deployments

## 🔄 Configuration Loading Strategy

### Intelligent Configuration Resolution

The system employs a robust configuration loading approach:

1. **Primary Configuration**:

   - Prioritizes unit-specific provider configurations
   - Dynamically loads provider settings from `unit_cfg_providers.hcl`
   - Applies version constraints from `unit_cfg_versions.hcl`

2. **Fallback Mechanism**:
   - Provides safe default configurations when specific files are missing
   - Includes a null provider to prevent initialization errors
   - Maintains infrastructure deployment capabilities

## 🛠️ Configuration Examples

### Basic Provider Setup

```hcl
# unit_cfg_providers.hcl
locals {
  providers = [
    <<-EOF
    provider "example" {
      # Provider-specific configuration
      credential = var.provider_credential
    }
    EOF
  ]
}

# unit_cfg_versions.hcl
locals {
  versions = [
    <<-EOF
    terraform {
      required_providers {
        example = {
          source  = "example/provider"
          version = "~> 1.0"
        }
      }
    }
    EOF
  ]
}
```

## 🔒 Security Considerations

- Never commit sensitive credentials in configuration files
- Use environment variables for credential management
- Follow the principle of least privilege
- Implement secure credential rotation strategies

## 🤝 Contributing Guidelines

When modifying the unit:

1. Update provider configurations in `unit_cfg_providers.hcl`
2. Modify version constraints in `unit_cfg_versions.hcl`
3. Validate changes using `terragrunt plan`
4. Ensure no sensitive information is exposed

## 🌐 Environment Variable Management

### Provider Credential Setup

```bash
# Generic provider credential example
export TG_STACK_PROVIDER_EXAMPLE_CREDENTIAL="your-secure-credential"
```

## 🔍 Troubleshooting

### Common Configuration Issues

1. **Provider Configuration Errors**

   - Verify environment variable names
   - Check credential formatting
   - Ensure correct provider source and version

2. **Version Constraint Problems**
   - Validate version syntax
   - Confirm provider compatibility
   - Check for version conflicts

## 📚 Related Documentation

- [Terragrunt Documentation](https://terragrunt.gruntwork.io/docs/)
- [Terraform Provider Development](https://developer.hashicorp.com/terraform/plugin/best-practices)

## Conclusion

This configuration approach provides a flexible, secure, and maintainable method for managing infrastructure providers across different deployment units.

## 🌳 Configuration Hierarchy and Inheritance

### Infrastructure Configuration Layers

The unit's configuration follows a multi-layered approach:

1. **Root Configuration** (`@root.hcl`):

   - Provides global infrastructure management logic
   - Implements dynamic provider and version generation
   - Manages shared configuration loading
   - Defines core Terragrunt generation rules

2. **Shared Configurations** (`@_shared/_config`):

   - Centralize common infrastructure metadata
   - Define reusable locals and configuration patterns
   - Provide baseline settings for remote state, tagging, and application metadata

3. **Unit Configuration** (`@terragrunt.hcl`):
   - Specific to this infrastructure unit
   - Inherits and extends root and shared configurations
   - Defines unit-specific resource modules and dependencies

### Configuration Flow

```
Root Config (@root.hcl)
│
├── Shared Configs (@_shared/_config)
│   ├── app.hcl
│   ├── remote_state.hcl
│   └── tags.hcl
│
└── Unit Config (@terragrunt.hcl)
    ├── unit_cfg_providers.hcl
    └── unit_cfg_versions.hcl
```

### Inheritance Mechanism

- **Provider Configuration**:

  - Root configuration dynamically loads provider settings from unit-level `unit_cfg_providers.hcl`
  - Fallback to null provider if no configuration is found

- **Version Management**:

  - Root configuration reads version constraints from `unit_cfg_versions.hcl`
  - Generates `versions.tf` with unit-specific or default constraints

- **Shared Metadata**:
  - Unit inherits common tags, remote state configuration, and application metadata
  - Allows for consistent resource management across infrastructure

### Configuration Precedence

1. Unit-specific configurations take highest priority
2. Shared configurations provide default values
3. Root configuration manages dynamic generation and fallback mechanisms

### Best Practices

- Keep unit configurations minimal and focused
- Leverage shared configurations for common settings
- Use environment variables for sensitive or environment-specific configurations
- Maintain clear separation of concerns between configuration layers
