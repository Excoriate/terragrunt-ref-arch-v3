# Disabled fixture: Sets the example's is_enabled flag to false.
# This should prevent the module call in main.tf, resulting in no
# data sources being read and null outputs from the module.

is_enabled = false
