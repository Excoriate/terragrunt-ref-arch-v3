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
