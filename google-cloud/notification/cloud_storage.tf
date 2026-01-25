module "cloud_storage" {
  source      = "../modules/cloud_storage"
  bucket_name = "budget-alert-to-slack"
}

# ソースコードをzipに圧縮
# Cloud Functionsはルートディレクトリにgo.modとエントリーポイントが必要
data "archive_file" "function_source" {
  type        = "zip"
  source_dir  = "${path.module}/src"
  output_path = "${path.module}/src/function_source.zip"

  # Cloud Functionsが期待しない不要なファイルを除外
  excludes = [
    "function_source.zip",
    ".env",
    "bin/**",
    "Taskfile.yaml",
    "internal/**/*_test.go",
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
