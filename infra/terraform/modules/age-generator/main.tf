resource "random_integer" "age" {
  min = local.validated_min_age
  max = local.validated_max_age
}
