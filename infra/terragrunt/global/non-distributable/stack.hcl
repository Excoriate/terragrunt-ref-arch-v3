// 🌐 Stack Configuration Manifest
//
// This file defines stack-specific configurations and metadata for the DNS infrastructure.
// It serves as a central point for stack-level settings that can be referenced across
// different Terragrunt and Terraform modules within the DNS domain.
//
// 🔍 Purpose:
// - Define stack-specific variables and metadata
// - Provide consistent tagging strategy for DNS-related resources
// - Enable stack-level customizations and grouping
//
// 💡 Configuration Guidelines:
// - Modify values to match specific DNS infrastructure requirements
// - Ensure consistency across different DNS-related components
// - Use meaningful, descriptive names and tags

locals {
  // 🏷️ Stack Naming Convention
  // - Use descriptive, lowercase names that represent the stack's primary function
  // - Recommended format: [domain-type-purpose]
  //
  // 💡 Tip: Stack name should clearly indicate its infrastructure domain and purpose
  stack_name = "non-distributable"

  // 🌐 Stack Description
  // Provides a human-readable description of the stack's purpose and scope
  stack_description = "implement non-distributable resources locally for testing purposes"

  // 📛 Stack-Level Tagging Strategy
  //
  // Tags provide crucial metadata for:
  // - Resource identification
  // - Logical grouping
  // - Infrastructure organization
  // - Compliance and management tracking
  //
  // 🔍 Best Practices:
  // - Extend environment-level tags with stack-specific metadata
  // - Use clear, descriptive tag values
  // - Consider additional context-specific tags
  tags = {
    // Stack-specific tags
    Stack   = "non-distributable"
    Domain  = "dev-utilities"
    Purpose = "implement non-distributable resources locally for testing purposes and local development"
  }

  // 🔧 Stack Configuration Flags
  // Enable or disable specific stack-wide configurations
  stack_config = {
  }
}
