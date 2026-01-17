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
