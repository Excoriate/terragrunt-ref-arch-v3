terraform {
  required_version = ">= 1.0" # Align with module

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0" # Align with module
    }
  }
}
