config {
  force = false
}

plugin "aws" {
  enabled = true
  version = "0.38.0" # Use consistent version
  source  = "github.com/terraform-linters/tflint-ruleset-aws"
}

plugin "terraform" {
  enabled = true
  preset  = "recommended"
}

# Mandatory Rules - Always Enabled
rule "terraform_deprecated_interpolation" {
  enabled = true
}

rule "terraform_deprecated_lookup" {
  enabled = true
}

rule "terraform_empty_list_equality" {
  enabled = true
}

rule "terraform_map_duplicate_keys" {
  enabled = true
}

rule "terraform_required_providers" {
  enabled = true
}

rule "terraform_required_version" {
  enabled = true
}

# Documentation Quality Rules
rule "terraform_documented_variables" {
  enabled = true
}

rule "terraform_documented_outputs" {
  enabled = true
}

# Code Quality Rules
rule "terraform_typed_variables" {
  enabled = true
}

rule "terraform_unused_declarations" {
  # Keep enabled for modules
  enabled = true
}

rule "terraform_unused_required_providers" {
  enabled = true
}
