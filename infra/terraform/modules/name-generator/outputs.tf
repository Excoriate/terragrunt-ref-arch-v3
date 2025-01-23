output "name" {
  description = "Randomly generated name"
  value       = random_shuffle.name.result[0]
}

output "full_name" {
  description = "Generated name with random suffix"
  value       = "${var.input_name}-${random_string.name_suffix.result}"
}

output "suffix" {
  description = "Generated random suffix"
  value       = random_string.name_suffix.result
}

output "gender" {
  description = "Gender of the generated name"
  value       = var.gender
}
