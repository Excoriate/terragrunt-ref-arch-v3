locals {
  # 📦 Versions configuration for Terragrunt
  # Specifies required provider version and source
  # TODO: Add version constraints, and uncomment the following block according to your unit's requirements
  versions = [
    <<-EOF
terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.0"
    }
  }
}
    EOF
  ]
}
