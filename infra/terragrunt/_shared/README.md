# Terragrunt Shared Infrastructure Configuration 🏗️

## Overview

This directory provides a centralized, modular approach to managing infrastructure configurations and reusable components for Terragrunt-based deployments.

## Architecture Principles 🛠️

### Core Design Concepts

1. **Centralization**: Single source of truth for infrastructure management
2. **Modularity**: Discrete, composable infrastructure components
3. **Flexibility**: Environment-aware, dynamically resolvable configurations
4. **Traceability**: Comprehensive metadata and resource identification

## Directory Structure 📂

```
_shared/
├── _config/       # Centralized configuration management
│   ├── README.md  # Detailed configuration strategy
│   ├── tags.hcl   # Resource tagging mechanism
│   ├── app.hcl    # Project metadata configuration
│   └── remote_state.hcl  # State management configuration
│
└── _units/        # Reusable infrastructure components
    ├── README.md  # Component architecture overview
    ├── dni_generator.hcl
    ├── lastname_generator.hcl
    ├── name_generator.hcl
    └── age_generator.hcl
```

## Configuration Management 🔧

### Shared Configuration (`_config/`)

Provides centralized, environment-aware configuration management:

- **Resource Tagging**: Implement consistent, hierarchical resource identification
- **Project Metadata**: Dynamic, environment-driven project configuration
- **Remote State Management**: Secure, flexible state storage strategies

**Detailed Documentation**: [\_config/README.md](_config/README.md)

### Infrastructure Units (`_units/`)

Offers modular, reusable infrastructure components:

- **Generator Components**: Discrete infrastructure building blocks
- **Dynamic Configuration**: Flexible, environment-specific customization
- **Intelligent Dependency Management**: Robust cross-module interactions

**Detailed Documentation**: [\_units/README.md](_units/README.md)

## Key Benefits 🌟

- **Consistent Configuration**: Uniform infrastructure management
- **Environment Flexibility**: Adaptable to different deployment contexts
- **Enhanced Traceability**: Comprehensive resource metadata
- **Simplified Maintenance**: Centralized, modular approach

## Getting Started 🚀

1. Review configuration strategies in `_config/README.md`
2. Explore available infrastructure units in `_units/README.md`
3. Set up required environment variables
4. Customize configurations for your specific infrastructure needs
