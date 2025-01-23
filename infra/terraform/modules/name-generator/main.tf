# Generate random suffix for the name
resource "random_string" "name_suffix" {
  length  = var.suffix_length
  special = false
  upper   = false
}

# Randomly select a name from the available names list
resource "random_shuffle" "name" {
  input        = local.available_names
  result_count = 1
}
