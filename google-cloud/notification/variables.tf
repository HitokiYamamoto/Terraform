variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
  sensitive   = true
}

variable "project_number" {
  description = "Google Cloud Project Number"
  type        = string
  sensitive   = true
}

variable "billing_account_id" {
  description = "Google Cloud Billing Account ID"
  type        = string
  sensitive   = true
}
