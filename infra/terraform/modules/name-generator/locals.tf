locals {
  # Predefined name lists with gender-specific variations
  male_names = ["Juan", "Carlos", "Miguel", "Pedro", "Luis"]
  female_names = ["Maria", "Ana", "Carmen", "Sofia", "Elena"]

  # Combine names based on gender
  available_names = var.gender == "male" ? local.male_names : var.gender == "female" ? local.female_names : concat(local.male_names, local.female_names)
}
