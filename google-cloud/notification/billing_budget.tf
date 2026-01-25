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
