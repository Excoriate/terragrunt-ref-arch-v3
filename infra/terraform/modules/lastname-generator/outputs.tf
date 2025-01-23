output "lastname" {
  description = "Randomly generated lastname"
  value       = random_shuffle.lastname.result[0]
}

output "full_lastname" {
  description = "Generated lastname with random suffix"
  value       = "${var.input_lastname}-${random_string.lastname_suffix.result}"
}

output "suffix" {
  description = "Generated random suffix"
  value       = random_string.lastname_suffix.result
}

output "gender" {
  description = "Gender of the generated lastname"
  value       = var.gender
}
