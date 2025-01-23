# 🌐 Global Default Terraform Variables
#
# Purpose: Provide baseline configuration defaults for infrastructure
#
# 🔍 Dynamic Variable Loading Mechanism:
# - Automatically loaded by Terragrunt during apply, plan, destroy commands
# - Can be overridden by environment or region-specific .tfvars files
# - Serves as a fallback configuration for infrastructure components
#
# 💡 Usage in root.hcl:
# - Dynamically included via extra_arguments "optional_vars"
# - Supports flexible, hierarchical configuration management
# - Enables environment-agnostic default settings
