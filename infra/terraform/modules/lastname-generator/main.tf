resource "random_string" "lastname_suffix" {
  length  = 4
  special = false
  upper   = false
}

resource "random_shuffle" "lastname" {
  input        = local.available_lastnames
  result_count = 1
}
