variable "repository_id" {
  description = "The ID of the GitHub repository where the branch protection will be applied."
  type        = string
}

variable "job_list" {
  description = "A list of GitHub Actions job names that must pass before merging."
  type        = list(string)
  default     = [] // 初回はCI設定がない場合を考慮して空リストをデフォルトにする
}
