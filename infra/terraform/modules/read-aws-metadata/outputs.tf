output "is_enabled" {
  description = "Indicates whether the data sources in this module were enabled (and therefore read)."
  value       = var.is_enabled
}

output "account_id" {
  description = "The AWS Account ID number of the account that owns or contains the calling entity."
  value       = var.is_enabled ? data.aws_caller_identity.current[0].account_id : null
}

output "caller_arn" {
  description = "The AWS ARN associated with the calling entity."
  value       = var.is_enabled ? data.aws_caller_identity.current[0].arn : null
}

output "caller_user_id" {
  description = "The unique identifier of the calling entity."
  value       = var.is_enabled ? data.aws_caller_identity.current[0].user_id : null
}

output "partition" {
  description = "The AWS partition in which the calling entity exists (e.g., `aws`, `aws-cn`, `aws-us-gov`)."
  value       = var.is_enabled ? data.aws_partition.current[0].partition : null
}

output "region_name" {
  description = "The name of the AWS Region (e.g., `us-east-1`)."
  value       = var.is_enabled ? data.aws_region.current[0].name : null
}

output "region_description" {
  description = "The description of the AWS Region (e.g., `US East (N. Virginia)`)."
  value       = var.is_enabled ? data.aws_region.current[0].description : null
}
