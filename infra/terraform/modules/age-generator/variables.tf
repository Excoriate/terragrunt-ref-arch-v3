variable "min_age" {
  description = "Minimum age for generation"
  type        = number
  default     = 18
  validation {
    condition     = var.min_age >= 0 && var.min_age <= 100
    error_message = "Minimum age must be between 0 and 100."
  }
}

variable "max_age" {
  description = "Maximum age for generation"
  type        = number
  default     = 65
  validation {
    condition     = var.max_age >= 18 && var.max_age <= 100
    error_message = "Maximum age must be between 18 and 100."
  }
}
