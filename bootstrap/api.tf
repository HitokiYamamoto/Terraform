locals {
  required_apis = [
    "artifactregistry.googleapis.com",     # Artifact Registry API
    "billingbudgets.googleapis.com",       # Billing Budgets API
    "cloudbilling.googleapis.com",         # Cloud Billing API
    "cloudbuild.googleapis.com",           # Cloud Build API
    "cloudfunctions.googleapis.com",       # Cloud Functions API
    "run.googleapis.com",                  # Cloud Run Admin API
    "cloudresourcemanager.googleapis.com", # Cloud Resource Manager API
    "eventarc.googleapis.com",             # Eventarc API
    "iamcredentials.googleapis.com",       # IAM Service Account Credentials API
    "iam.googleapis.com",                  # Identity and Access Management (IAM) API
    "pubsub.googleapis.com",               # Pub/Sub API
    "secretmanager.googleapis.com",        # Secret Manager API
    "sts.googleapis.com",                  # Security Token Service (STS) API
  ]
}

resource "google_project_service" "required_apis" {
  for_each = toset(local.required_apis)

  project = var.project_id
  service = each.key

  disable_on_destroy = false
}
