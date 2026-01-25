module "service_account" {
  source       = "../modules/service_account"
  account_id   = "budget-alert-function"
  display_name = "Budget Alert Function Service Account"
}

resource "google_project_iam_member" "service_account_build_roles" {
  for_each = toset(local.build_roles)

  project = var.project_id
  role    = each.value
  member  = module.service_account.member
}

resource "google_secret_manager_secret_iam_binding" "secret_accessor" {
  for_each  = local.secrets
  secret_id = each.value
  role      = "roles/secretmanager.secretAccessor"
  members = [
    module.service_account.member
  ]
}

resource "google_cloud_run_service_iam_binding" "invoker" {
  location = google_cloudfunctions2_function.function.location
  service  = google_cloudfunctions2_function.function.name
  role     = "roles/run.invoker"
  members = [
    module.service_account.member,
  ]
}
