# 🌐 Root Terragrunt Configuration Inclusion
# This block imports the common configuration from the parent directory's terragrunt.hcl file.
# It enables consistent configuration sharing across multiple Terragrunt modules, ensuring that
# all modules can access shared settings and parameters defined at the root level.
include "root" {
  path           = find_in_parent_folders("root.hcl")
  merge_strategy = "deep"
}

# 🧩 Shared Units Configuration
# This block imports standardized component configuration from the shared components directory.
# It is important to note that any modifications should be made in the shared component configuration file
# located at: `_shared/_components/quota-generator.hcl`. This ensures that changes are reflected
# across all modules that utilize this shared configuration.
include "shared" {
  path           = "${get_terragrunt_dir()}/../../../_shared/_units/name_generator.hcl"
  expose         = true
  merge_strategy = "deep"
}

locals {
  # 🔧 Terraform Module Source Resolution
  # ---------------------------------------------------------------------------------------------------------------------
  # Intelligent Terraform module source management with:
  # - Local path testing support
  # - Version override capability
  # - Fallback to default version
  # - Dynamic source path computation
  #
  # Resolution Strategy:
  # 1. If local_path is provided, use it for local testing
  # 2. If version_override is set, use that version
  # 3. Otherwise, fall back to the default version from shared configuration
  #
  # Enables flexible module sourcing for:
  # - Local development and testing
  # - Precise version control
  # - Consistent module referencing across infrastructure
  # ---------------------------------------------------------------------------------------------------------------------
  tf_module_local_path       = "${include.shared.locals.git_base_url}/name-generator"
  tf_module_version_override = ""
  tf_module_version          = local.tf_module_version_override != "" ? local.tf_module_version_override : include.shared.locals.tf_module_version_default
  tf_module_source           = include.shared.locals.tf_module_source

  tg_source_computed    = local.tf_module_local_path != "" ? local.tf_module_local_path : format("%s?ref=%s", local.tf_module_source, local.tf_module_version)
  echo_tf_module_source = run_cmd("sh", "-c", "echo '🔧  TF Module Source (child): ${local.tg_source_computed}'")
}

# 🚀 Terraform Source Configuration
# Dynamically resolves the Terraform module source based on:
# - Local testing path (if provided)
# - Shared module source
# - Specific version reference
terraform {
  source = local.tg_source_computed
}

# 📦 Inputs Configuration
# Provides an empty inputs block to satisfy Terragrunt configuration requirements
# Actual inputs are managed in the shared configuration
inputs = {}
