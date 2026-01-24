locals {
  service_account_roles = [
    "roles/storage.admin",
    "roles/iam.serviceAccountAdmin",
    "roles/iam.serviceAccountUser", # サービスアカウントの利用権限
    "roles/artifactregistry.admin",
    "roles/cloudfunctions.admin",
    "roles/resourcemanager.projectIamAdmin", # プロジェクトレベルのIAMポリシーを変更する権限
  ]
}

# ユーザーアカウントにサービスアカウントのimpersonate権限を付与
resource "google_service_account_iam_binding" "token_creator" {
  service_account_id = module.default_service_account.id
  role               = "roles/iam.serviceAccountTokenCreator"

  members = var.impersonators
}

# サービスアカウントに各種管理権限を付与
resource "google_project_iam_member" "service_account_roles" {
  for_each = toset(local.service_account_roles)

  project = var.project_id
  role    = each.value
  member  = module.default_service_account.member
}
