config {
  force = false
}

plugin "aws" {
  enabled = true
  version = "0.38.0" # Use a consistent version across examples/modules
  source  = "github.com/terraform-linters/tflint-ruleset-aws"
}

plugin "terraform" {
  enabled = true
  preset  = "recommended"
}

# Recommended Rules (subset relevant for examples)
rule "terraform_deprecated_index" {
  enabled = true
}

rule "terraform_deprecated_interpolation" {
  enabled = true
}

rule "terraform_documented_outputs" {
  enabled = true
}

rule "terraform_documented_variables" {
  enabled = true
}

rule "terraform_naming_convention" {
  enabled = true
}

rule "terraform_required_providers" {
  enabled = true
}

rule "terraform_required_version" {
  enabled = true
}

rule "terraform_typed_variables" {
  enabled = true
}

rule "terraform_unused_declarations" {
  enabled = true
}

# Rules typically disabled for examples
rule "terraform_standard_module_structure" {
  enabled = false # Examples are root modules
}

rule "terraform_module_pinned_source" {
  enabled = false # Examples call local modules via relative path
}

rule "terraform_module_version" {
  enabled = false # Examples call local modules
}
