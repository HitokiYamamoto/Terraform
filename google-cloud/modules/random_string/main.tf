resource "random_string" "main" {
  length  = 6
  lower   = true
  upper   = false
  special = false
  numeric = true
}
