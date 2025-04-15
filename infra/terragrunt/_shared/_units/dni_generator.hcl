locals {
  # ---------------------------------------------------------------------------------------------------------------------
  # 🏗️ ENVIRONMENT CONFIGURATION ORCHESTRATION
  # ---------------------------------------------------------------------------------------------------------------------
  # Purpose: Dynamically load and aggregate environment-specific configuration files
  # This mechanism enables a flexible, layered infrastructure configuration approach
  #
  # Configuration Layers:
  # - Environment-level settings (env.hcl)
  # - Global architecture configurations
  # - Hierarchical tag management
  #
  # Key Benefits:
  # - Modular configuration management
  # - Centralized environment settings
  # - Flexible tag inheritance
  # ---------------------------------------------------------------------------------------------------------------------
  env_cfg = read_terragrunt_config(find_in_parent_folders("env.hcl"))

  # 🌐 Stack Configuration
  # Loads the stack configuration file that serves as the single source of truth
  # for stack-level settings and metadata
  stack_cfg = read_terragrunt_config(find_in_parent_folders("stack.hcl"))

  # ---------------------------------------------------------------------------------------------------------------------
  # 🌐 GLOBAL ARCHITECTURE CONFIGURATION RESOLVER
  # ---------------------------------------------------------------------------------------------------------------------
  # Purpose: Centralize and standardize infrastructure-wide configuration
  # Loads the root configuration file that serves as the single source of truth
  # for infrastructure-level settings and metadata
  #
  # Key Responsibilities:
  # - Provide global configuration context
  # - Enable consistent infrastructure metadata
  # - Support cross-module configuration sharing
  # ---------------------------------------------------------------------------------------------------------------------
  cfg = read_terragrunt_config("${find_in_parent_folders("config.hcl")}")

  # ---------------------------------------------------------------------------------------------------------------------
  # 🏷️ INTELLIGENT TAG ORCHESTRATION SYSTEM
  # ---------------------------------------------------------------------------------------------------------------------
  # Purpose: Create a sophisticated, hierarchical tag management mechanism
  # Implements a multi-layered tagging strategy that allows for:
  # - Global tag inheritance
  # - Environment-specific tag augmentation
  # - Unit-level tag customization
  #
  # Tag Hierarchy (from broadest to most specific):
  # 1. Global Architecture Tags 🌐
  # 2. Environment-Level Tags 🌍
  # 3. Unit-Specific Tags 🧩
  # 4. Stack-Level Tags 📚
  #
  # Benefits:
  # - Consistent resource identification
  # - Flexible tag management
  # - Enhanced resource tracking and compliance
  # ---------------------------------------------------------------------------------------------------------------------
  unit_tags = {
    Unit = "dni-generator"
    Type = "random-generator"
  }

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

  # ---------------------------------------------------------------------------------------------------------------------
  # 🔧 GIT MODULE SOURCE CONFIGURATION
  # ---------------------------------------------------------------------------------------------------------------------
  # Intelligent Terraform module source management with:
  # - Centralized default configuration
  # - Flexible repository and path selection
  # - Semantic version control
  #
  git_base_url              = local.cfg.locals.cfg_git.git_base_urls.github
  tf_module_repository      = "your-org/terraform-modules.git"
  tf_module_version_default = get_env("TG_STACK_TF_MODULE_DNI_GENERATOR_VERSION_DEFAULT", "v0.1.0")
  tf_module_path_default    = "modules/dni-generator"

  tf_module_source = format(
    "%s%s//%s",
    local.git_base_url,
    local.tf_module_repository,
    local.tf_module_path_default
  )

  echo_tf_module_source = run_cmd("sh", "-c", "echo '🔧  TF Module Source (parent): ${local.tf_module_source} version: ${local.tf_module_version_default}'")

  # ---------------------------------------------------------------------------------------------------------------------
  # 🌐 SPECIFIC MODULE CONFIGURATION
  # ---------------------------------------------------------------------------------------------------------------------
  # Here we define the specific configuration for the module.
  # This is useful when we need to override the default configuration for the module.
  #
  # TODO: Add specific locals, or configuration for your module to be computed here.
}

# 🔗 DEPENDENCIES
# These blocks define dependencies for the Terragrunt configuration.
# Dependencies allow for the management of resources that rely on
# other configurations, ensuring that they are created or updated
# in the correct order. This promotes modularity and reusability
# of infrastructure components.

dependency "age_generator" {
  config_path = "${find_in_parent_folders("dni/age_generator")}"
  mock_outputs = {
    generated_age = 30
  }
}

dependency "name_generator" {
  config_path = "${find_in_parent_folders("dni/name_generator")}"
  mock_outputs = {
    full_name = "john-abc123"
  }
}

dependency "lastname_generator" {
  config_path = "${find_in_parent_folders("dni/lastname_generator")}"
  mock_outputs = {
    full_lastname = "smith-xyz789"
  }
}

dependencies {
  paths = [
    "${find_in_parent_folders("dni/age_generator")}",
    "${find_in_parent_folders("dni/name_generator")}",
    "${find_in_parent_folders("dni/lastname_generator")}"
  ]
}

# 🚀 TERRAGRUNT INFRASTRUCTURE UNIT CONFIGURATION
# Defines the input parameters for the infrastructure unit
# Combines global configuration, metadata, and tag management
inputs = {
  # Ensure all variables from variables.tf are explicitly set
  prefix                  = ""   # Default empty string
  generate_control_letter = true # Default true
  name                    = dependency.name_generator.outputs.full_name
  lastname                = dependency.lastname_generator.outputs.full_lastname

  # CRITICAL: Explicitly set the age variable
  # Use the dependency output or provide a default/mock value
  age = dependency.age_generator.outputs.generated_age

  country = "Spain" # Default country
  tags    = local.all_tags
}
