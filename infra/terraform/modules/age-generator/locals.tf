locals {
  # No complex local computations needed at this time
  # Ensure min_age is less than max_age
  validated_min_age = var.min_age < var.max_age ? var.min_age : var.max_age
  validated_max_age = var.max_age > var.min_age ? var.max_age : var.min_age
}
