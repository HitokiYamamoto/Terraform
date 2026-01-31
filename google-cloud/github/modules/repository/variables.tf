variable "repository_name" {
  description = "GitHubのリポジトリ名"
  type        = string
}

variable "description" {
  description = "リポジトリの説明"
  type        = string
  default     = ""
}

variable "has_issues" {
  description = "Issue機能を有効にするか"
  type        = bool
  default     = true
}

variable "has_wiki" {
  description = "Wiki機能を有効にするか"
  type        = bool
}

variable "auto_init" {
  description = "READMEを自動作成するか"
  type        = bool
}

variable "topics" {
  description = "リポジトリのトピック（タグ）"
  type        = list(string)

  validation {
    condition     = alltrue([for t in var.topics : length(t) <= 50 && can(regex("^[a-z0-9]+(-[a-z0-9]+)*$", t))])
    error_message = "各トピックは50文字以下の小文字英数字およびハイフンで構成され、ハイフンで始まってはいけません。"
  }
}

variable "delete_branch_on_merge" {
  description = "マージ後にブランチを自動削除するか"
  type        = bool
}

variable "visibility" {
  description = "パブリックリポジトリ or プライベートリポジトリ"
  type        = string

  validation {
    condition     = contains(["public", "private"], var.visibility)
    error_message = "visibilityはpublic, privateのいずれかを指定してください。"
  }
}

variable "archived" {
  description = "リポジトリをアーカイブするか"
  type        = bool
}
