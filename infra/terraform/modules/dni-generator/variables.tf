variable "prefix" {
  description = "Optional prefix for the DNI number (first digits)"
  type        = string
  default     = ""
  validation {
    condition     = can(regex("^\\d{0,8}$", var.prefix))
    error_message = "Prefix must be a string of digits with a maximum length of 8."
  }
}

variable "generate_control_letter" {
  description = "Whether to generate the control letter"
  type        = bool
  default     = true
}

variable "name" {
  description = "First name for DNI generation"
  type        = string
}

variable "lastname" {
  description = "Last name for DNI generation"
  type        = string
}

variable "age" {
  description = "Age for DNI generation"
  type        = number
}

variable "country" {
  description = "Country for DNI generation"
  type        = string
  default     = "Spain"
  validation {
    condition     = contains(["Spain", "Argentina", "Mexico"], var.country)
    error_message = "Country must be 'Spain', 'Argentina', or 'Mexico'."
  }
}

variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}
