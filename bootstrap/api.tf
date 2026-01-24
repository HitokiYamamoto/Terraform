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
/* Identity and Access Management (IAM) API
# サービスアカウントやIAMポリシーを管理するためのAPI
*/
resource "google_project_service" "iam" {
  project = var.project_id
  service = "iam.googleapis.com"

  disable_on_destroy = false
}

/* Artifact Registry API
# Dockerコンテナイメージなどのアーティファクトを管理するためのAPI
*/
resource "google_project_service" "artifactregistry" {
  project = var.project_id
  service = "artifactregistry.googleapis.com"

  disable_on_destroy = false
}

/* Cloud Functions API
# Cloud Functions (第2世代) を作成・管理するためのAPI
*/
resource "google_project_service" "cloudfunctions" {
  project = var.project_id
  service = "cloudfunctions.googleapis.com"

  disable_on_destroy = false
}

/* Cloud Build API
# Cloud Functionsのデプロイ時にコンテナイメージをビルドするためのAPI
*/
resource "google_project_service" "cloudbuild" {
  project = var.project_id
  service = "cloudbuild.googleapis.com"

  disable_on_destroy = false
}

/* Cloud Run Admin API
# Cloud Functions (第2世代) の実行基盤となるCloud Runを管理するためのAPI
*/
resource "google_project_service" "run" {
  project = var.project_id
  service = "run.googleapis.com"

  disable_on_destroy = false
}

/* Cloud Resource Manager API
# プロジェクト情報の取得やIAMポリシーの管理に必要なAPI
# google_project_iam_memberなどのリソースで使用される
*/
resource "google_project_service" "cloudresourcemanager" {
  project = var.project_id
  service = "cloudresourcemanager.googleapis.com"

  disable_on_destroy = false
}

/* Pub/Sub API
# Pub/Subトピックの作成・管理に必要なAPI
*/
resource "google_project_service" "pubsub" {
  project = var.project_id
  service = "pubsub.googleapis.com"

  disable_on_destroy = false
}

/* Cloud Billing API
# 予算アラートの作成・管理に必要なAPI
*/
resource "google_project_service" "cloud_billing" {
  project = var.project_id
  service = "cloudbilling.googleapis.com"

  disable_on_destroy = false
}

/* Billing Budgets API
# 予算アラートの作成・管理に必要なAPI
*/
resource "google_project_service" "billing_budgets" {
  project = var.project_id
  service = "billingbudgets.googleapis.com"

  disable_on_destroy = false
}

/* Eventarc API
# Cloud Functions (第2世代) のイベントトリガーに必要なAPI
*/
resource "google_project_service" "eventarc" {
  project            = var.project_id
  service            = "eventarc.googleapis.com"
  disable_on_destroy = false
}
