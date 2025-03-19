// 🌐 Terragrunt Global Configuration
// Centralized infrastructure configuration management system that orchestrates shared settings,
// provider management, and global variables across the entire infrastructure ecosystem.
// This configuration serves as the single source of truth for infrastructure-wide settings.

locals {
  // 📂 Shared Configuration Loading
  // Standardizes and centralizes configuration across infrastructure units by reading
  // configuration files from a predefined shared directory. This approach enables modular
  // and reusable infrastructure settings, ensuring consistency and reducing duplication.
  // When adding new shared configurations, maintain the established directory structure.
  # shared_config_app = read_terragrunt_config("${get_parent_terragrunt_dir()}/_shared/_config/app.hcl")
  shared_config_app          = read_terragrunt_config("${get_terragrunt_dir()}/_shared/_config/app.hcl")
  shared_config_remote_state = read_terragrunt_config("${get_terragrunt_dir()}/_shared/_config/remote_state.hcl")
  shared_config_tags         = read_terragrunt_config("${get_terragrunt_dir()}/_shared/_config/tags.hcl")

  // 🔌 Provider Configuration Management
  // Safely read provider configurations if they exist
  unit_cfg_providers = try(
    read_terragrunt_config("${get_original_terragrunt_dir()}/unit_cfg_providers.hcl"),
    {
      locals = {
        providers = []
      }
    }
  )

  unit_cfg_versions = try(
    read_terragrunt_config("${get_original_terragrunt_dir()}/unit_cfg_versions.hcl"),
    {
      locals = {
        versions = []
      }
    }
  )

  // 🛡️ Provider Configuration Override
  // Flag to control provider file generation
  // When set to "true", forces provider file generation even if no providers are defined
  // When set to "false", prevents provider file generation if no providers are found
  // Default behavior (null/unset) allows generation based on configuration presence
  enable_providers_override = get_env("TG_STACK_FLAG_ENABLE_PROVIDERS_OVERRIDE", "true")

  // 🌐 Dynamic Provider Configuration
  shared_config_providers = local.unit_cfg_providers.locals

  // 🔌 Provider Configuration Management
  // Flag to control version file generation
  // When set to "true", forces version file generation even if no versions are defined
  // When set to "false", prevents version file generation if no versions are found
  // Default behavior (null/unset) allows generation based on configuration presence
  enable_versions_override = get_env("TG_STACK_FLAG_ENABLE_VERSIONS_OVERRIDE", "true")
  generate_versions_file   = length(local.unit_cfg_versions.locals.versions) > 0 && (local.enable_versions_override == "true" || local.enable_versions_override == null) ? true : false
  generate_providers_file  = length(local.shared_config_providers.providers) > 0 && (local.enable_providers_override == "true" || local.enable_providers_override == null) ? true : false

  // 🔍 Configuration Extraction
  // Simplifies access to nested configuration structures by creating direct references
  // to local blocks in shared configuration files. This approach reduces complex nested
  // access and provides a clean, straightforward way to retrieve configuration values.
  // Helps maintain readability and reduces the chance of configuration access errors.
  cfg_app          = local.shared_config_app.locals
  cfg_remote_state = local.shared_config_remote_state.locals
  cfg_tags         = local.shared_config_tags.locals

  // 🏷️ Global Resource Tagging
  // Implements consistent resource identification and management through a centralized
  // tagging strategy. These tags enable cost tracking, compliance monitoring, and
  // comprehensive resource management across the entire infrastructure.
  // Ensure all tags are meaningful, descriptive, and follow organizational standards.
  tags = local.cfg_tags.tags

  // 🔖 Project Metadata Management
  // Maintains consistent project identification across infrastructure by providing
  // a single source of truth for project-level metadata. This ensures that project
  // name and version are consistently applied throughout all infrastructure components.
  // Update these values in the shared app configuration file to propagate changes.
  product_name    = local.cfg_app.product_name
  product_version = local.cfg_app.product_version

  // 🌍 Runtime Environment Configuration
  // Provides flexible, environment-driven configuration with sensible defaults to
  // support multi-environment deployments. Allows runtime configuration flexibility
  // through environment variables, enabling easy environment-specific customizations.
  // Use environment variables to override default settings when needed.
  deployment_region_unnormalized = get_env("TG_STACK_REGION", "us-east-1")
  deployment_region              = lower(trimspace(local.deployment_region_unnormalized))

  // 💾 Remote State Management
  // Configures and standardizes Terraform state storage by centralizing remote state
  // configuration. Ensures consistent state management across all infrastructure units,
  // providing a reliable and predictable approach to tracking infrastructure changes.
  // Modify remote state settings in the remote_state.hcl configuration file as your
  // infrastructure requirements evolve.
  bucket_name           = local.cfg_remote_state.bucket_name
  lock_table            = local.cfg_remote_state.lock_table
  region                = local.cfg_remote_state.region
  state_object_basename = local.cfg_remote_state.state_object_basename
  backend_tf_filename   = local.cfg_remote_state.backend_tf_filename

  // 📦 Terraform Version Management
  // Centralized configuration for Terraform version control and file generation.
  // These settings control both the enforced version and the generation of .terraform-version files.
  terraform_version_disabled_comments = <<-EOT
# Terraform Version File Generation Disabled
# To enable version file generation, set TG_STACK_FLAG_ENABLE_TERRAFORM_VERSION_FILE_OVERRIDE to "true"
# And, set TG_STACK_TF_VERSION to the desired Terraform version
# Current setting prevents automatic version file creation
    EOT

  // 🔒 Version Enforcement Configuration
  // To modify the enforced Terraform version:
  // 1. Set TG_STACK_TF_VERSION environment variable, or
  // 2. Update the default version here (currently "1.9.0")
  tf_version_enforced_unnormalised = get_env("TG_STACK_TF_VERSION", "1.9.0")
  tf_version_enforced              = lower(trimspace(local.tf_version_enforced_unnormalised))

  // 🎛️ Version File Generation Control
  // To enable .terraform-version file generation:
  // 1. Set TG_STACK_FLAG_ENABLE_TERRAFORM_VERSION_FILE_OVERRIDE to "true"
  // 2. Default is "false" to prevent automatic generation
  // When enabled: Creates .terraform-version with tf_version_enforced
  // When disabled: Creates file with instructions (see terraform_version_disabled_comments)
  enable_terraform_version_file_override = get_env("TG_STACK_FLAG_ENABLE_TERRAFORM_VERSION_FILE_OVERRIDE", "false")
}
