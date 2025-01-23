output "generated_age" {
  description = "Generated age within specified range"
  value       = random_integer.age.result
}

output "min_age" {
  description = "Minimum age used for generation"
  value       = local.validated_min_age
}

output "max_age" {
  description = "Maximum age used for generation"
  value       = local.validated_max_age
}
