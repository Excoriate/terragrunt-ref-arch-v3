variable "aws_region" {
  description = "AWS region for the example."
  type        = string
  default     = "us-east-1"
}

variable "is_enabled" {
  description = "Flag to enable/disable the module call in this example."
  type        = bool
  default     = true
}
