resource "random_string" "dni_prefix" {
  length  = 1
  special = false
  upper   = false
  numeric = true
}

resource "random_shuffle" "dni_prefix" {
  input        = local.available_dni_prefixes
  result_count = 1
}

resource "random_integer" "dni_number" {
  min = 10000000
  max = 99999999
}
