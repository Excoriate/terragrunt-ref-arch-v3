variable "is_enabled" {
  type        = bool
  description = <<-DESC
  Controls whether the random string resource is created.
  Set to `false` to prevent the random string from being generated.

  **Purpose**: Allows conditional creation of the random string.
  **Impact**: If `false`, the `random_string` output will be `null`.
  **Default**: `true` (module is enabled by default)
  DESC
  default     = true
}

variable "length" {
  type        = number
  description = <<-DESC
  Specifies the exact length of the random string to be generated.

  **Purpose**: Controls the length of the output random string.
  **Impact**: Directly influences the length of the `random_string` output.
  **Default**: `8`
  **Constraints**: Must be a positive integer. The `random` provider might have practical upper limits.
  DESC
  default     = 8
  validation {
    condition     = var.length > 0
    error_message = "Length (var.length) must be a positive number."
  }
}

variable "lower" {
  type        = bool
  description = <<-DESC
  Specifies whether lowercase letters are allowed in the generated random string.

  **Purpose**: Controls the character set used for generation.
  **Impact**: Affects the possible characters in the `random_string` output.
  **Default**: `true` (lowercase letters are included by default)
  DESC
  default     = true
}

variable "upper" {
  type        = bool
  description = <<-DESC
  Specifies whether uppercase letters are allowed in the generated random string.

  **Purpose**: Controls the character set used for generation.
  **Impact**: Affects the possible characters in the `random_string` output.
  **Default**: `true` (uppercase letters are included by default)
  DESC
  default     = true
}
