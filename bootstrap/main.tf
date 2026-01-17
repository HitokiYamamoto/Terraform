/* NOTE:
このファイルにあるコードは、TerraformのBootstrap問題を解決するためのもの。
以下のリソースを手動作成や、ユーザー権限で作成しているので取扱注意

1. Terraform実行用のサービスアカウント
2. Backend用のGCSバケット
*/

/* IAM Service Account Credentials API
# サービスアカウントの認証情報（アクセストークンなど）を生成・管理するためのAPI
# OIDC認証後にGoogle Cloudリソースへアクセスする際に必要
*/
resource "google_project_service" "iamcredentials" {
  project = var.project_id
  service = "iamcredentials.googleapis.com"

  disable_on_destroy = false
}

/* Security Token Service (STS) API
# 外部IDプロバイダー（GitHub）のOIDCトークンをGoogle Cloudの認証情報に交換するAPI
# Workload Identity連携の核となるサービス
*/
resource "google_project_service" "sts" {
  project = var.project_id
  service = "sts.googleapis.com"

  disable_on_destroy = false
}

module "default_service_account" {
  source       = "../google-cloud/modules/service_account"
  account_id   = "default-terraform-sa"
  display_name = "Terraform実行用のサービスアカウント"
}

module "backend_bucket" {
  source      = "../google-cloud/modules/cloud_storage"
  bucket_name = "terraform-backend-bucket"
}

resource "google_storage_bucket_iam_binding" "admin" {
  bucket = module.backend_bucket.bucket_name
  role   = "roles/storage.admin"

  members = [
    module.default_service_account.member,
  ]
}

module "github_oidc" {
  source         = "../google-cloud/modules/github_oidc"
  project_id     = var.project_id
  github_org_id  = var.github_org_id
  github_repo_id = var.github_repo_id
}

resource "google_service_account_iam_binding" "workload_identity_user" {
  service_account_id = module.default_service_account.id
  role               = "roles/iam.workloadIdentityUser"

  members = [
    module.github_oidc.principal,
  ]
}
