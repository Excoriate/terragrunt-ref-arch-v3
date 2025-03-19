locals {
  // 🔧 Configuration File Loading
  // Centralizes infrastructure configuration by reading global and environment-specific config files.
  // This mechanism allows for modular and flexible infrastructure management across different environments.
  // When adding new configuration sources, ensure they follow the established path and naming conventions.
  # cfg = read_terragrunt_config("${get_parent_terragrunt_dir()}/config.hcl")
  cfg       = read_terragrunt_config("${find_in_parent_folders("config.hcl")}")
  env_cfg   = read_terragrunt_config("${get_terragrunt_dir()}/../../env.hcl")
  stack_cfg = read_terragrunt_config("${get_terragrunt_dir()}/../stack.hcl")

  // 🌍 Deployment Context Extraction
  // Captures critical runtime environment details like deployment region and environment name.
  // These variables enable dynamic configuration and support multi-environment infrastructure strategies.
  // Modify environment-specific configurations in the respective env.hcl files.
  deployment_region        = local.cfg.locals.deployment_region
  deployment_environment   = local.env_cfg.locals.environment_name
  deployment_stack         = local.stack_cfg.locals.stack_name
  path_relative_to_include = path_relative_to_include()

  // 🎯 Current Infrastructure Unit Identification
  // Dynamically determines the specific infrastructure unit being processed using Terragrunt's path resolution.
  // This approach allows for unit-specific configurations and operations without hardcoding unit names.
  // Ensure your directory structure maintains the expected hierarchy for accurate unit identification.
  current_unit = basename(local.path_relative_to_include)

  // 🗂️ Deployment Group and State Management
  // Generates a consistent, predictable state folder structure that supports complex infrastructure hierarchies.
  // Creates unique identifiers by combining product name and deployment group with a flat naming convention.
  // When restructuring infrastructure, maintain the logical relationship between product name and deployment paths.
  deployment_group = replace(local.path_relative_to_include, "/", "-")
  terraform_state_folder_flat = replace(
    format("%s-%s",
      local.cfg.locals.product_name,
      local.deployment_group
    ),
    "/",
    "-"
  )

  // 🔑 Remote State Key Path Generation
  // Constructs a standardized, hierarchical key for Terraform remote state that reflects the infrastructure's logical structure.
  // Format: <product_name>/<environment>/<unit-path>-<state_object_basename>
  // Example: myapp/global/dns-dns_zone-terraform.tfstate.json
  remote_state_key_path = join("/", [
    local.cfg.locals.product_name,
    local.deployment_environment,
    format("%s-%s",
      replace(trimprefix(local.path_relative_to_include, "${local.deployment_environment}/"), "/", "-"),
      local.cfg.locals.state_object_basename
    )
  ])

  // 🔌 Dynamic Provider Loading
  // Enables flexible, unit-specific provider configuration by dynamically loading provider settings.
  // Supports multi-provider setups and allows granular control over provider registration for different infrastructure units.
  // To add a new provider, update the shared provider configuration files and register the unit appropriately.
  dynamic_providers = try(
    local.cfg.locals.unit_cfg_providers.locals.providers != null
    ? local.cfg.locals.unit_cfg_providers.locals.providers
    : [],
    []
  )

  // 📦 Dynamic Versions Configuration
  // Manages Terraform provider version constraints dynamically, ensuring compatibility across different infrastructure units.
  // Loads version configurations specific to registered units, preventing version conflicts.
  // When introducing new providers or updating versions, modify the corresponding provider configuration files.
  dynamic_versions = try(
    local.cfg.locals.unit_cfg_versions.locals.versions != null
    ? local.cfg.locals.unit_cfg_versions.locals.versions
    : [],
    []
  )

  // 🛡️ Fallback Provider Configuration
  // Ensures every infrastructure unit has a valid provider configuration by supplying a default 'null' provider.
  // Prevents configuration errors when no unit-specific providers are defined.
  // Serves as a safety mechanism to maintain infrastructure deployment capabilities.
  fallback_providers = [
    <<-EOF
# Default provider for units without specific provider configuration
provider "null" {}
    EOF
  ]

  // 🏷️ Fallback Versions Configuration
  // Provides a minimal Terraform versions configuration as a last resort for infrastructure units.
  // Prevents initialization errors by guaranteeing a valid versions block.
  // Can be extended to include more default provider version constraints if needed.
  fallback_versions = [
    <<-EOF
# Default versions block for units without specific version configuration
terraform {
  required_providers {
    null = {
      source = "hashicorp/null"
      version = "~> 3.2"
    }
  }
}
    EOF
  ]

  // 🔬 Provider and Versions Resolution
  // Implements an intelligent selection mechanism for provider and version configurations.
  // Prioritizes unit-specific dynamic configurations while maintaining a robust fallback strategy.
  // Ensures flexibility and reliability in infrastructure provisioning.
  final_providers = length(local.dynamic_providers) > 0 ? local.dynamic_providers : local.fallback_providers

  final_versions = length(local.dynamic_versions) > 0 ? local.dynamic_versions : local.fallback_versions

  // 📢 Deployment Information Logging
  // Provides visibility into the current Terragrunt execution context by outputting critical deployment details.
  // Helps in debugging and tracking infrastructure changes during execution.
  // Can be extended to include additional diagnostic information if needed.
  echo_line_separator                      = run_cmd("sh", "-c", "echo '================================================================================'")
  echo_current_unit                        = run_cmd("sh", "-c", "echo '🏗️  Current Unit: ${local.current_unit}'")
  echo_product_name                        = run_cmd("sh", "-c", "echo '📦  Product Name: ${local.cfg.locals.product_name}'")
  echo_product_version                     = run_cmd("sh", "-c", "echo '🔧  Product Version: ${local.cfg.locals.product_version}'")
  echo_deployment_group                    = run_cmd("sh", "-c", "echo '📦  Deployment Group: ${local.deployment_group}'")
  echo_deployment_region                   = run_cmd("sh", "-c", "echo '🌍  Deployment Region: ${local.deployment_region}'")
  echo_deployment_stack                    = run_cmd("sh", "-c", "echo '🏞️  Deployment Stack: ${local.deployment_stack}'")
  echo_environment_name                    = run_cmd("sh", "-c", "echo '🏞️  Environment Name: ${local.deployment_environment}'")
  echo_remote_state_key_path               = run_cmd("sh", "-c", "echo '🔑  Remote State Key Path: ${local.remote_state_key_path}'")
  echo_is_using_providers_then_true        = run_cmd("sh", "-c", "echo '🔌  Is Using Providers: ${length(local.dynamic_providers) > 0 ? "true" : "false"}'")
  echo_is_overwriting_versionstf_then_true = run_cmd("sh", "-c", "echo '🔧  Is Overwriting Versions: ${length(local.dynamic_versions) > 0 ? "true" : "false"}'")

  // 📢 State Information
  echo_state_bucket_name = run_cmd("sh", "-c", "echo '🗄️  State Bucket Name: ${local.cfg.locals.bucket_name}'")
  echo_state_lock_table  = run_cmd("sh", "-c", "echo '🔒  State Lock Table: ${local.cfg.locals.lock_table}'")
}

terraform {
  // 🔧 Flexible Variable File Handling
  // Terragrunt makes it easy to manage different configuration files across environments.
  // We can seamlessly include global, environment-specific, and region-specific variable files.
  // This approach lets you customize your infrastructure settings without duplicating code.
  extra_arguments "optional_vars" {
    commands = [
      "apply",
      "destroy",
      "plan",
    ]

    optional_var_files = [
      // 📁 Global Default Variables
      // Start with a base configuration that applies everywhere
      "default.tfvars",

      // 🌍 Environment-Specific Tweaks
      // Override or extend global settings for specific environments
      "${local.deployment_environment}/default.tfvars",

      // 🗺️ Region-Specific Customizations
      // Fine-tune configurations for different regions
      "${local.deployment_region}.tfvars",
    ]
  }
}

// 💾 Smart Remote State Management
// Keep your Terraform state secure, organized, and easily accessible.
// We use S3 as a centralized backend with built-in encryption and tagging.
remote_state {
  backend = "s3"
  generate = {
    path      = local.cfg.locals.backend_tf_filename
    if_exists = "overwrite"
  }

  config = {
    disable_bucket_update = true
    encrypt               = true

    region         = local.cfg.locals.region
    bucket         = local.cfg.locals.bucket_name
    dynamodb_table = local.cfg.locals.lock_table

    key = local.remote_state_key_path

    s3_bucket_tags      = local.cfg.locals.tags
    dynamodb_table_tags = local.cfg.locals.tags
  }
}

// 🚀 Dynamic Provider Configuration
// This block generates provider settings dynamically for each infrastructure unit.
// By doing so, it allows for a more flexible and manageable configuration of providers,
// adapting to the specific needs of each unit without hardcoding values.
// The 'path' specifies where the generated provider configuration will be saved.
// The 'if_exists' condition checks if the generation of the provider file is enabled;
// if so, it will overwrite any existing file, otherwise it will skip the generation.
// The 'contents' field constructs the provider configuration by joining the list of
// shared provider configurations, ensuring that only relevant settings are included.
generate "providers" {
  path      = "providers.tf"
  if_exists = local.cfg.locals.generate_providers_file ? "overwrite" : "skip"
  contents  = local.cfg.locals.generate_providers_file ? join("\n", local.cfg.locals.shared_config_providers.providers) : ""
}

// 🏗️ Intelligent Version Management
// Dynamically handle provider versions across different infrastructure units.
// Ensures compatibility and makes version updates a breeze.
generate "versions" {
  path      = "versions.tf"
  if_exists = local.cfg.locals.generate_versions_file ? "overwrite" : "skip"
  contents  = local.cfg.locals.generate_versions_file ? join("\n", local.cfg.locals.unit_cfg_versions.locals.versions) : ""
}

// 🔧 Terraform Version Enforcement
// Automatically create a .terraform-version file to keep everyone on the same page.
// Control which Terraform version is used across your infrastructure.
generate "terraform_version" {
  path      = ".terraform-version"
  if_exists = "overwrite_terragrunt"
  contents  = local.cfg.locals.enable_terraform_version_file_override == "true" ? local.cfg.locals.tf_version_enforced : local.cfg.locals.terraform_version_disabled_comments
}
