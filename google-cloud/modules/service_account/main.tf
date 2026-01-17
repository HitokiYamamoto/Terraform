resource "google_service_account" "main" {
  account_id                   = var.account_id
  display_name                 = var.display_name
  disabled                     = false
  create_ignore_already_exists = true // 既に存在する場合は作成を無視
}
