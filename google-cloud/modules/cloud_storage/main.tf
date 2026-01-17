module "random_string" {
  source = "../random_string"
}

resource "google_storage_bucket" "main" {
  name = "${var.bucket_name}-${module.random_string.result}"
  # ※ "US" (マルチリージョン) は無料枠の対象外なので注意
  location                    = "us-central1" // 無料枠を考慮してUSリージョンを指定
  storage_class               = "STANDARD"    // 一番安い
  force_destroy               = false
  uniform_bucket_level_access = true // ファイルごとの権限設定を禁止。IAMだけで管理する。

  versioning {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true // 誤って削除するのを防止
  }

  lifecycle_rule {
    action {
      type = "Delete"
    }

    condition {
      num_newer_versions = 3
    }
  }
}
