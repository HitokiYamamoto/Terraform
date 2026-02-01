module "terraform_protect" {
  source        = "./modules/branch_protection"
  repository_id = module.terraform.node_id
  job_list = [
    "ci-result",           # Terraformワークフローの統合チェック
    "check-golang-result", # Goファイルチェック
  ]
}
