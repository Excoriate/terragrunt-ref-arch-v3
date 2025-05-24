data "aws_caller_identity" "current" {
  count = var.is_enabled ? 1 : 0
}

data "aws_partition" "current" {
  count = var.is_enabled ? 1 : 0
}

data "aws_region" "current" {
  count = var.is_enabled ? 1 : 0
}
