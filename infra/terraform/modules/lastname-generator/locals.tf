locals {
  # Predefined lastname lists with gender-specific variations
  male_lastnames = ["Garcia", "Rodriguez", "Martinez", "Lopez", "Gonzalez"]
  female_lastnames = ["Garcia", "Rodriguez", "Martinez", "Lopez", "Gonzalez"]

  # Combine lastnames based on gender
  available_lastnames = var.gender == "male" ? local.male_lastnames : var.gender == "female" ? local.female_lastnames : concat(local.male_lastnames, local.female_lastnames)
}
