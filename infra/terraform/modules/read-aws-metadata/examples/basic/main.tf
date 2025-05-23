module "this" {
  source = "../../" # Reference the parent module directory

  is_enabled = var.is_enabled # Controlled by fixtures
}
