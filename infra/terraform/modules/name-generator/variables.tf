variable "input_name" {
  description = "Base name to concatenate with random string"
  type        = string
  validation {
    condition     = length(var.input_name) > 0
    error_message = "Input name must not be empty."
  }
}

variable "suffix_length" {
  description = "Length of the random suffix"
  type        = number
  default     = 6
  validation {
    condition     = var.suffix_length > 0 && var.suffix_length <= 16
    error_message = "Suffix length must be between 1 and 16 characters."
  }
}

variable "gender" {
  description = "Gender for name generation"
  type        = string
  default     = "any"
  validation {
    condition     = contains(["male", "female", "any"], var.gender)
    error_message = "Gender must be 'male', 'female', or 'any'."
  }
}
