terraform {
  required_version = ">= 1.11.3"

  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6.0"
    }
  }
}
