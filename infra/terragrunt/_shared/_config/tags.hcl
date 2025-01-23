# ---------------------------------------------------------------------------------------------------------------------
# 🌍 COMMON TAGS CONFIGURATION
# This block establishes a standardized set of metadata tags applicable to all infrastructure resources
# managed by Terraform and orchestrated by Terragrunt within the project. Tags are key-value pairs associated
# with resources that serve multiple purposes, including identification, organization, and governance of
# resources across cloud environments. 🏗️

# Utilizing a consistent tagging strategy across all modules enhances the ability to:
# - 🔍 Identify resources, their purpose, and their lifecycle owner at a glance.
# - 💰 Implement cost allocation, reporting, and optimization strategies based on tags.
# - 🔒 Enforce security and compliance policies through tag-based resource segmentation.
# - ⚙️ Automate operations and management tasks that depend on the categorization of resources.

# This block defines a reusable set of global tags that can be incorporated into any module by including it in the
# `locals` block. Subsequently, these tags can be merged with module-specific tags and applied to resources in the
# Terraform `resource` blocks, ensuring a unified and comprehensive tagging approach across the project's infrastructure.
# ---------------------------------------------------------------------------------------------------------------------
locals {
  tags = {
    ManagedBy      = "Terraform"                          // 🛠️ Indicates the tool used for resource provisioning.
    OrchestratedBy = "Terragrunt"                         // 🎛️ Indicates the tool used for workflow orchestration.
    Author         = get_env("TG_STACK_APP_AUTHOR", "") // ✍️ The author of the configuration.
    Type           = "infrastructure"                     // 📦 Categorizes the resource within the broader infrastructure ecosystem.
    Application    = get_env("TG_STACK_APP_PRODUCT_NAME", "my-app") // 📱 The application or service that the resource supports.
    # TODO: Add git-sha tag. Uncomment when ready. Ensure there's at least one commit in the repository.
    # "git-sha"      = run_cmd("--terragrunt-global-cache", "--terragrunt-quiet", "git", "rev-parse", "--short", "HEAD")
  }
}
