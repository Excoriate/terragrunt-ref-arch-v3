# Environment Architecture Framework 🌐

## Overview

This directory represents a standardized approach to environment configuration in infrastructure-as-code, providing a flexible, modular framework for managing complex infrastructure deployments across different contexts and environments.

## Architectural Principles 🏗️

### Core Design Concepts

1. **Hierarchical Configuration**: Multi-level configuration management
2. **Environment Abstraction**: Consistent, context-independent infrastructure definition
3. **Dynamic Configuration**: Adaptive, context-aware settings
4. **Metadata Enrichment**: Comprehensive resource identification and tracking

## Reference Architecture Structure 📂

```
environment/
├── env.hcl           # Environment-level configuration manifest
├── .envrc            # Environment variable management
├── default.tfvars    # Baseline variable configurations
└── stack-name/       # Infrastructure stack
    ├── stack.hcl     # Stack-specific configuration
    └── units/        # Modular infrastructure components
        ├── unit-a/
        ├── unit-b/
        └── unit-c/
```

## Configuration Files Deep Dive 🔍

### 1. Environment Configuration Manifest (`env.hcl`) 🏷️

#### Purpose

A centralized configuration file that defines environment-specific settings, metadata, and tagging strategies.

#### Key Components

- **Environment Naming**: Standardized identification
- **Tagging Strategy**: Consistent resource metadata
- **Naming Conventions**: Structured resource identification

#### Example Configuration

```hcl
locals {
  # Environment Identification
  environment_name = "dev"  # dev, staging, prod, global
  environment      = "dev"  # Short identifier

  # Resource Tagging Strategy
  tags = {
    Environment = local.environment_name
    ManagedBy   = "Terragrunt"
    Project     = "Infrastructure"
  }
}
```

#### Best Practices

- Use lowercase, descriptive environment names
- Maintain consistent tagging across all resources
- Include metadata that aids in resource management

### 2. Environment Variable Management (`.envrc`) 🌍

#### Purpose

A bash script that provides robust environment variable management, logging, and configuration detection.

#### Key Features

- **Dynamic .env File Loading**: Searches for and loads `.env` files
- **Environment Root Detection**: Identifies infrastructure configuration roots
- **Secure Variable Export**: Safely exports and logs environment variables
- **Configuration Validation**: Ensures critical variables are defined

#### Core Functions

- `_env_init()`: Initializes environment configuration
- `_load_env_dotenv()`: Loads environment variables from `.env` files
- `_env_export()`: Securely exports environment-specific variables
- `_validate_env_config()`: Validates critical configuration settings

#### Example Usage

```bash
# Automatically sets environment variables
# Logs configuration details
# Provides flexible, secure configuration management
```

#### Best Practices

- Never commit sensitive information
- Use environment-specific `.env` files
- Implement least-privilege access controls

### 3. Default Terraform Variables (`default.tfvars`) 📋

#### Purpose

Provides baseline configuration defaults for infrastructure components.

#### Key Characteristics

- Automatically loaded by Terragrunt
- Serves as a fallback configuration
- Can be overridden by environment-specific `.tfvars` files

#### Configuration Strategy

- Define default values for infrastructure components
- Enable flexible, hierarchical configuration management
- Support environment-agnostic default settings
