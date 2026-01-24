module "budget_alert" {
  source             = "./notification"
  project_id         = var.project_id
  project_number     = var.project_number
  billing_account_id = var.billing_account_id
}
