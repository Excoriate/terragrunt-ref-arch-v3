variable "is_enabled" {
  description = "Flag to enable/disable the execution of data sources within this module."
  type        = bool
  default     = true
}

variable "tags" {
  description = "A map of tags to assign to resources. Although this module creates no taggable resources, this variable is included for consistency."
  type        = map(string)
  default     = {}
}
