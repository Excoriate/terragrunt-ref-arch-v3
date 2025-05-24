locals {
  # 🌐 PROVIDER CREDENTIAL MANAGEMENT
  # Purpose: Securely handle and normalize provider credentials
  #
  # Key Features:
  # - Environment variable-based credential retrieval
  # - Automatic normalization (lowercase, trimmed)
  # - Flexible error handling for missing credentials

  # 🔑 Raw Credential Retrieval
  # Captures original, unmodified environment variable values
  # TODO: Add environment variable retrieval, or specific credential retrieval for your provider
  # provider_credential_unnormalized = get_env("TG_STACK_PROVIDER_CREDENTIAL", "")

  # 🧼 Credential Normalization
  # Applies consistent formatting to credentials
  # - Converts to lowercase
  # - Removes leading/trailing whitespaces
  # - Handles empty input gracefully
  # TODO: Add normalization logic, and uncomment the following block
  # provider_credential = local.provider_credential_unnormalized != "" ? lower(trimspace(local.provider_credential_unnormalized)) : ""

  # 📋 Provider Configuration
  # Centralizes provider-specific settings and credentials
  # provider_config = {
  #   credential = local.provider_credential
  # }

  # 🎲 RANDOM PROVIDER CONFIGURATION
  # Purpose: Configure the Random provider for generating random values
  #
  # Key Features:
  # - No credentials required
  # - Simple provider configuration
  # - Used for generating random values in a deterministic way

  # ⚙️ Provider configuration for Terragrunt
  # Generates the provider block with normalized credentials
  providers = [
    <<-EOF
provider "random" {
}
    EOF
  ]
}
