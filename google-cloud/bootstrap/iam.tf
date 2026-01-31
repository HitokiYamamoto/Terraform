locals {
  service_account_roles = [
    "roles/storage.admin",
    "roles/iam.serviceAccountAdmin",
    "roles/artifactregistry.admin",
    "roles/cloudfunctions.admin",
    "roles/resourcemanager.projectIamAdmin", # プロジェクトレベルのIAMポリシーを変更する権限
    "roles/pubsub.admin",
    "roles/secretmanager.admin",
  ]
  # dataリソースで取得したいサービスアカウント名のリスト
  target_service_account_names = [
    "budget-alert-function"
  ]
}

# 特定のサービスアカウントをdataリソースで取得
data "google_service_account" "targets" {
  for_each   = toset(local.target_service_account_names)
  account_id = each.value
  project    = var.project_id
}

# ユーザーアカウントにサービスアカウントのimpersonate権限を付与
resource "google_service_account_iam_binding" "token_creator" {
  service_account_id = module.default_service_account.id
  role               = "roles/iam.serviceAccountTokenCreator"

  members = var.impersonators
}

# 特定のサービスアカウントに対してのみroles/iam.serviceAccountUserを付与
resource "google_service_account_iam_binding" "service_account_user" {
  for_each = data.google_service_account.targets

  service_account_id = each.value.id
  role               = "roles/iam.serviceAccountUser"

  members = [
    module.default_service_account.member,
  ]
}

# サービスアカウントに各種管理権限を付与
resource "google_project_iam_member" "service_account_roles" {
  for_each = toset(local.service_account_roles)

  project = var.project_id
  role    = each.value
  member  = module.default_service_account.member
}

# サービスアカウントに請求管理者権限を付与
resource "google_billing_account_iam_member" "billing_admin" {
  billing_account_id = var.billing_account_id
  role               = "roles/billing.admin"
  member             = module.default_service_account.member
}
