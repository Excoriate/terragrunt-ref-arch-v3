output "dni_prefix" {
  description = "Randomly generated DNI prefix"
  value       = random_shuffle.dni_prefix.result[0]
}

output "dni_number" {
  description = "Randomly generated DNI number"
  value       = random_integer.dni_number.result
}

output "full_dni" {
  description = "Complete DNI with prefix and number"
  value       = "${random_shuffle.dni_prefix.result[0]}${random_integer.dni_number.result}"
}

output "country" {
  description = "Country of DNI generation"
  value       = var.country
}

output "full_name" {
  description = "Full name of the generated citizen"
  value       = local.full_name
}

output "age" {
  description = "Age of the generated citizen"
  value       = var.age
}
