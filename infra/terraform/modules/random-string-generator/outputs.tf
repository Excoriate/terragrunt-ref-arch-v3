output "random_string" {
  description = "The generated random string based on the specified length, lower, and upper character constraints. Returns 'null' if 'is_enabled' is false."
  value       = var.is_enabled ? random_string.this[0].result : null
}
