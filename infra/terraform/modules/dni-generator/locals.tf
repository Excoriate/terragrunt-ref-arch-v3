locals {
  # DNI control letter calculation
  dni_control_letters = ["T", "R", "W", "A", "G", "M", "Y", "F", "P", "D", "X", "B", "N", "J", "Z", "S", "Q", "V", "H", "L", "C", "K", "E"]

  # Generate a deterministic prefix based on name, lastname, and age
  name_hash = substr(md5("${var.name}${var.lastname}"), 0, 4)
  age_hash = substr(md5("${var.age}"), 0, 4)
  prefix = substr("${local.name_hash}${local.age_hash}", 0, 6)

  base_dni_number = format("%08d", random_integer.dni_number.result)
  control_letter = local.dni_control_letters[tonumber(local.base_dni_number) % 23]
  full_dni = "${local.base_dni_number}${local.control_letter}"
  full_name = "${var.name} ${var.lastname}"

  # DNI generation logic based on country
  spain_dni_prefix = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9"]
  argentina_dni_prefix = ["1", "2", "3", "4", "5", "6", "7", "8", "9"]
  mexico_dni_prefix = ["1", "2", "3", "4", "5", "6", "7", "8", "9"]

  # Select DNI prefix based on country
  available_dni_prefixes = var.country == "Spain" ? local.spain_dni_prefix : var.country == "Argentina" ? local.argentina_dni_prefix : local.mexico_dni_prefix
}
