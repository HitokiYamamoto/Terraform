module "suffix" {
  source = "../random_string"
}

resource "google_iam_workload_identity_pool" "main" {
  description               = "GitHub からOIDCでGoogle Cloudに接続するためのworkload identity pool"
  disabled                  = false
  display_name              = "GitHub OIDC Pool"
  project                   = var.project_id
  workload_identity_pool_id = "github-oidc-pool-${module.suffix.result}"
}

resource "google_iam_workload_identity_pool_provider" "main" {
  project                            = var.project_id
  workload_identity_pool_id          = google_iam_workload_identity_pool.main.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-oidc-provider-${module.suffix.result}"
  display_name                       = "GitHub OIDC Provider"
  description                        = "GitHub Actions用のOIDC Provider"

  attribute_mapping = {
    "google.subject"                = "assertion.sub"
    "attribute.repository_id"       = "assertion.repository_id"
    "attribute.repository_owner_id" = "assertion.repository_owner_id"
  }

  attribute_condition = <<EOT
    assertion.repository_owner_id == "${var.github_org_id}" &&
    attribute.repository_id == "${var.github_repo_id}"
  EOT

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}
