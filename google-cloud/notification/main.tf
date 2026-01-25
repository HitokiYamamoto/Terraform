locals {
  build_roles = [
    "roles/logging.logWriter",
    "roles/artifactregistry.writer",
    "roles/storage.objectViewer",
    "roles/cloudbuild.builds.builder",
  ]

  secrets = ({
    SLACK_BOT_USER_OAUTH_TOKEN = module.slack_oauth_token.secret_id
    CHANNEL_NAME               = module.slack_channel_name.secret_id,
  })
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
  source_dir  = "${path.module}/src"
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

module "slack_oauth_token" {
  source    = "../modules/secret_manager"
  secret_id = "slack-bot-user-oauth-token"
}

module "slack_channel_name" {
  source    = "../modules/secret_manager"
  secret_id = "slack-channel-name"
}

# roles/secretmanager.secretAccessor
resource "google_secret_manager_secret_iam_binding" "secret_accessor" {
  for_each  = local.secrets
  secret_id = each.value
  role      = "roles/secretmanager.secretAccessor"
  members = [
    module.service_account.member
  ]
}

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

resource "google_cloud_run_service_iam_binding" "invoker" {
  location = google_cloudfunctions2_function.function.location
  service  = google_cloudfunctions2_function.function.name
  role     = "roles/run.invoker"
  members = [
    module.service_account.member,
  ]
}

/* Pub/Subの作成
Cloud Functions（特に第2世代）において Pub/Sub トリガーを利用する場合、
Cloud Functions（Eventarc）に管理を任せたほうが、構築・運用コストが圧倒的に低いので、
サブスクリプションを作ると逆に面倒になるため、Topicのみを作成する(と、Geminiが言ってた...)
*/
# メインのトピック
resource "google_pubsub_topic" "budget_alert_topic" {
  name    = "budget-alert-topic"
  project = var.project_id
}

# デッドレター用トピック（エラーになったメッセージの墓場）
resource "google_pubsub_topic" "budget_alert_dead_letter" {
  name    = "budget-alert-dead-letter"
  project = var.project_id
}

# デッドレター用トピックの権限（Pub/Subがここに書き込めるようにする）
resource "google_pubsub_topic_iam_binding" "dlq_publisher" {
  topic = google_pubsub_topic.budget_alert_dead_letter.name
  role  = "roles/pubsub.publisher"
  members = [
    "serviceAccount:service-${var.project_number}@gcp-sa-pubsub.iam.gserviceaccount.com"
  ]
}

# サブスクリプションの権限（Pub/SubがここからAckできるようにする）
resource "google_pubsub_subscription_iam_binding" "dlq_subscriber" {
  subscription = google_pubsub_subscription.budget_alert_subscription.name
  role         = "roles/pubsub.subscriber"
  members = [
    "serviceAccount:service-${var.project_number}@gcp-sa-pubsub.iam.gserviceaccount.com"
  ]
}

# メインのサブスクリプション
resource "google_pubsub_subscription" "budget_alert_subscription" {
  name  = "budget-alert-sub"
  topic = google_pubsub_topic.budget_alert_topic.name

  # Push設定：Cloud Functions (Cloud Run) のURLを叩く
  push_config {
    push_endpoint = google_cloudfunctions2_function.function.service_config[0].uri

    # セキュリティ設定：このSAの権限を使って関数を叩く
    oidc_token {
      service_account_email = module.service_account.email
    }
  }

  # デッドレターポリシー
  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.budget_alert_dead_letter.id
    max_delivery_attempts = 5 # 5回失敗したらデッドレターへ送る
  }

  # 再試行ポリシー（即時再送による攻撃を防ぐ）
  retry_policy {
    minimum_backoff = "10s"  # 最初は10秒待つ
    maximum_backoff = "600s" # 最大10分待つ
  }
}

# デッドレタートピック用のサブスクリプション（メッセージ保管庫）
resource "google_pubsub_subscription" "budget_alert_dead_letter_sub" {
  name  = "budget-alert-dead-letter-sub"
  topic = google_pubsub_topic.budget_alert_dead_letter.name

  # Pull型にする（push_configを書かなければ自動的にPull型になる）

  # メッセージの保持期間（最大7日）
  # これを設定しておかないと、調査する前に消える可能性がある
  message_retention_duration = "604800s" # 7日間

  # サブスクリプション自体の有効期限（無期限にする設定）
  # 31日間アクティブでないと削除されるのを防ぐ
  expiration_policy {
    ttl = ""
  }
}

# 予算アラートの作成
resource "google_billing_budget" "budget_any_cost" {
  billing_account = var.billing_account_id
  display_name    = "【緊急】課金発生アラート"

  budget_filter {
    projects = ["projects/${var.project_number}"]
  }

  amount {
    specified_amount {
      currency_code = "JPY"
      units         = "1" # 予算を「1円」に設定
    }
  }

  threshold_rules {
    threshold_percent = 1.0 # 1円の100% = 1円を超えたら通知
  }

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.budget_alert_topic.id
  }
}

resource "google_billing_budget" "budget_standard" {
  billing_account = var.billing_account_id
  display_name    = "月次予算管理 (1万円)"

  budget_filter {
    projects = ["projects/${var.project_number}"]
  }

  amount {
    specified_amount {
      currency_code = "JPY"
      units         = "10000"
    }
  }

  # 25%(2000円), 50%(5000円), 90%(9000円), 100%(1万円) で通知
  threshold_rules {
    threshold_percent = 0.25
  }
  threshold_rules {
    threshold_percent = 0.5
  }
  threshold_rules {
    threshold_percent = 0.9
  }
  threshold_rules {
    threshold_percent = 1.0
  }

  all_updates_rule {
    pubsub_topic = google_pubsub_topic.budget_alert_topic.id
  }
}
