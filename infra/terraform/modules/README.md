# Terraform Modules

## Overview

This directory contains a collection of modular Terraform components designed to provide a flexible, reusable, and scalable infrastructure-as-code (IaC) approach within our Terragrunt-based reference architecture.

## Architecture Philosophy

The modules in this directory embody key principles of modern infrastructure design:

- **Modularity**: Each module represents a discrete, self-contained infrastructure component
- **Reusability**: Modules are crafted to be environment-agnostic and easily composable
- **Flexibility**: Supports multiple sourcing strategies for enhanced development and deployment workflows

## Module Sourcing Strategies

The architecture supports multiple module sourcing mechanisms, as demonstrated in the `terragrunt.hcl` configuration:

```hcl
locals {
  tf_module_local_path       = "${get_repo_root()}/infra/terraform/modules/dni-generator"
  tf_module_version_override = ""
  tf_module_version          = local.tf_module_version_override != "" ? local.tf_module_version_override : include.shared.locals.tf_module_version_default
  tf_module_source           = include.shared.locals.tf_module_source
}

terraform {
  source = local.tf_module_local_path != "" ? local.tf_module_local_path : format("%s?ref=%s", local.tf_module_source, local.tf_module_version)
}
```

### Sourcing Mechanisms

1. **Local Development Path**

   - During development, modules can be sourced directly from the local filesystem
   - Enables rapid iteration and testing without pushing changes to a remote repository
   - Set `tf_module_local_path` to the local module directory

2. **Version-Controlled Remote Source**

   - Modules can be sourced from a remote repository with specific version references
   - Supports consistent, reproducible infrastructure deployments
   - Version controlled through `tf_module_version`

3. **Fallback Mechanism**
   - Intelligent fallback to default module source if no local path is specified
   - Ensures flexibility across different development and deployment environments

## Module Structure

Each module typically contains:

- `main.tf`: Primary resource definitions
- `variables.tf`: Input variable declarations
- `outputs.tf`: Module output definitions
- `locals.tf`: Local value computations
- `versions.tf`: Provider and Terraform version constraints

## Best Practices

- **Minimal Complexity**: Each module should have a single, well-defined responsibility
- **Parameterization**: Maximize configurability through input variables
- **Consistent Naming**: Use clear, descriptive names that reflect the module's purpose
- **Documentation**: Maintain comprehensive inline documentation

## Example Module Usage

```hcl
module "example_generator" {
  source = "path/to/module"

  # Module-specific input variables
  name        = var.name
  environment = var.environment
}
```

## Integration with Terragrunt

Modules are seamlessly integrated with Terragrunt through:

- Shared configuration files
- Dynamic source resolution
- Environment-specific parameter injection
