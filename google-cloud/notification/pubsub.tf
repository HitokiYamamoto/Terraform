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

# メインのサブスクリプション
resource "google_pubsub_subscription" "budget_alert_subscription" {
  name  = "budget-alert-sub"
  topic = google_pubsub_topic.budget_alert_topic.name

  # Push設定：Cloud Functions (Cloud Run) のURLを叩く
  push_config {
    push_endpoint = "${google_cloudfunctions2_function.function.service_config[0].uri}/projects/${var.project_id}/topics/${google_pubsub_topic.budget_alert_topic.name}"

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

# --- サービスエージェントへの権限付与 --- #

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
