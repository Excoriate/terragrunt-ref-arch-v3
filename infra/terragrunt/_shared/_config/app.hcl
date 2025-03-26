// 🎯 Local Variables for Terragrunt Configuration
// This section defines local variables that will be used across different child configurations.
// These locals help maintain consistency in naming conventions and metadata throughout the architecture.
locals {
  // 📦 Product Name
  product_name = get_env("TG_STACK_APP_PRODUCT_NAME", "my-app") // The name of the project for identification purposes.
  // 📦 Product Version
  product_version = "0.0.1" // The version of the project.
}
