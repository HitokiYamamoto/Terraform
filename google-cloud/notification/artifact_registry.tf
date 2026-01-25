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
