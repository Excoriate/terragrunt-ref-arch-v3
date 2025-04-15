resource "random_string" "this" {
  count = var.is_enabled ? 1 : 0

  length  = var.length
  lower   = var.lower
  upper   = var.upper
  special = false
  numeric = false
}
