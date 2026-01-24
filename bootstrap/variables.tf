variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
  sensitive   = true
}

variable "github_org_id" {
  description = "GitHub Organization ID"
  type        = string
  sensitive   = true
}

variable "github_repo_id" {
  description = "GitHub Repository ID"
  type        = string
  sensitive   = true
}

variable "impersonators" {
  description = "The member to grant the impersonate role to"
  type        = list(string)
  sensitive   = true
}

variable "billing_account_id" {
  description = "Google Cloud Billing Account ID"
  type        = string
  sensitive   = true
}
