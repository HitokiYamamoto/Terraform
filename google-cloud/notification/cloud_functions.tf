resource "google_cloudfunctions2_function" "function" {
  name        = "budget-alert-function"
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
    ingress_settings      = "ALLOW_ALL" # セキュリティ強化のため内部トラフィックのみ許可

    dynamic "secret_environment_variables" {
      for_each = local.secrets
      content {
        key        = secret_environment_variables.key
        secret     = secret_environment_variables.value
        project_id = var.project_id
        version    = "latest"
      }
    }
  }
}
