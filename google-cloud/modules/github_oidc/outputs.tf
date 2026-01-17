output "principal" {
  description = "OIDC Principal"
  value       = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.main.name}/attribute.repository_id/${var.github_repo_id}"
}
