locals {
  build_roles = [
    "roles/logging.logWriter",
    "roles/artifactregistry.writer",
    "roles/storage.objectViewer",
  ]
}

# サービスアカウントにビルド時に必要な権限を付与
resource "google_project_iam_member" "service_account_build_roles" {
  for_each = toset(local.build_roles)

  project = var.project_id
  role    = each.value
  member  = module.service_account.member
}

module "service_account" {
  source       = "../modules/service_account"
  account_id   = "budget-alert-function"
  display_name = "Budget Alert Function Service Account"
}

module "cloud_storage" {
  source      = "../modules/cloud_storage"
  bucket_name = "budget-alert-to-slack"
}

# ソースコードをzipに圧縮
# Cloud Functionsはルートディレクトリにgo.modとエントリーポイントが必要
data "archive_file" "function_source" {
  type        = "zip"
  source_dir  = path.module
  output_path = "${path.module}/function_source.zip"

  # Cloud Functionsが期待しない不要なファイルを除外
  excludes = [
    "function_source.zip",
    ".env",
    "bin/**",
    "Taskfile.yaml",
    "*.tf",
    "*.tfvars",
    ".terraform/**",
    "cmd/**", # ローカル開発用のcmdディレクトリは除外
  ]
}

# zipファイルをCloud Storageにアップロード
resource "google_storage_bucket_object" "function_source" {
  name   = "function_source-${data.archive_file.function_source.output_md5}.zip"
  bucket = module.cloud_storage.bucket_name
  source = data.archive_file.function_source.output_path

  depends_on = [
    data.archive_file.function_source
  ]
}

resource "google_artifact_registry_repository" "function_repo" {
  location      = "us-central1"
  repository_id = "slack-budget-alert-functions"
  description   = "Repository for Cloud Functions with cleanup policy"
  format        = "DOCKER"

  # 30日より古いイメージを削除
  cleanup_policies {
    id     = "delete-old-images"
    action = "DELETE"
    condition {
      older_than = "2592000s" # 30日
    }
  }

  # 最新の3つは残す
  # DELETEルールよりもKEEPルールが優先される
  cleanup_policies {
    id     = "keep-recent-versions"
    action = "KEEP"
    most_recent_versions {
      keep_count            = 3
      package_name_prefixes = [] # すべてのパッケージ対象
    }
  }
}

resource "google_cloudfunctions2_function" "function" {
  name        = "gcf-function"
  location    = "us-central1"
  description = "Budget Alert to Slack Function"

  build_config {
    runtime         = "go124"
    entry_point     = "ProcessBudgetAlert" # 実際のエントリーポイント関数名
    service_account = module.service_account.id
    /* NOTE:
    何も設定しない場合、裏側でソースコードをビルドしてコンテナイメージにし、Artifact Registryに保存されるらしい。
    月間0.5GB までの無料枠があるが、ビルドのたびにイメージが保存されていくため、明示的にArtifact Registryを作成し、古いイメージを削除するポリシーを設定する。
    */
    docker_repository = google_artifact_registry_repository.function_repo.id

    source {
      storage_source {
        bucket = module.cloud_storage.bucket_name
        object = google_storage_bucket_object.function_source.name
      }
    }
  }

  service_config {
    max_instance_count    = 1 # 跳ね上がり防止でmaxを設定
    min_instance_count    = 0 # コスト削減のためにminを0に設定
    available_memory      = "256M"
    timeout_seconds       = 60
    service_account_email = module.service_account.email
  }
}
