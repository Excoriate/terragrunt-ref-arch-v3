output "is_enabled" {
  description = "Output from the module: Indicates whether the data sources were enabled."
  value       = module.this.is_enabled
}

output "account_id" {
  description = "Output from the module: The AWS Account ID."
  value       = module.this.account_id
}

output "caller_arn" {
  description = "Output from the module: The AWS ARN associated with the calling entity."
  value       = module.this.caller_arn
}

output "caller_user_id" {
  description = "Output from the module: The unique identifier of the calling entity."
  value       = module.this.caller_user_id
}

output "partition" {
  description = "Output from the module: The AWS partition."
  value       = module.this.partition
}

output "region_name" {
  description = "Output from the module: The AWS Region name."
  value       = module.this.region_name
}

output "region_description" {
  description = "Output from the module: The AWS Region description."
  value       = module.this.region_description
}
