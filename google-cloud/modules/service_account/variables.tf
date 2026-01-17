variable "account_id" {
  description = "The ID of the service account."
  type        = string

  /* NOTE:
  6~30文字、英小文字、数字、ハイフン、アンダースコアのみ許可
  個人的にハイフンが好きなので、アンダースコアは禁止にする
  */
  validation {
    condition     = length(var.account_id) >= 6 && length(var.account_id) <= 30 && can(regex("^[a-z0-9-]+$", var.account_id))
    error_message = "The account_id must be 6-30 characters long and can only contain lowercase letters, numbers, and hyphens."
  }
}

variable "display_name" {
  description = "The display name of the service account."
  type        = string
}
