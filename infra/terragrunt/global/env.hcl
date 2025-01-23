// 🌍 Environment Configuration Manifest
//
// This file defines environment-specific configurations and metadata.
// It serves as a central point for environment-level settings that can be
// referenced across different Terragrunt and Terraform modules.
//
// 🔍 Purpose:
// - Define environment-specific variables
// - Provide consistent tagging strategy
// - Enable environment-level customizations
//
// 💡 Configuration Guidelines:
// - Modify values to match your specific environment requirements
// - Ensure consistency across different infrastructure components
// - Use meaningful, descriptive names for environments

// 🌍 Define local variables for the environment
locals {
  // 🏷️ Environment Naming Convention
  // - Use descriptive, lowercase names
  // - Recommended formats:
  //   * Development: "dev"
  //   * Staging: "staging"
  //   * Production: "prod"
  //   * Global/Shared: "global"
  //
  // 💡 Tip: Consistent naming helps in resource identification and management
  environment_name = "global"

  // 🌐 Short Environment Identifier
  // - Useful for resource naming, tagging, and quick reference
  // - Should match the full environment name or be a clear abbreviation
  environment = "global"

  // 📛 Resource Tagging Strategy
  //
  // Tags provide crucial metadata for:
  // - Resource identification
  // - Cost allocation
  // - Access management
  // - Compliance tracking
  //
  // 🔍 Best Practices:
  // - Keep tags consistent across all resources
  // - Use clear, descriptive tag values
  // - Consider adding more tags like:
  //   * Project
  //   * ManagedBy
  //   * Owner
  //   * CostCenter
  tags = {
    // Primary environment identifier
    // TODO: Add more tags
    Environment = local.environment_name
  }
}
