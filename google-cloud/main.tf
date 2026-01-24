module "budget_alert" {
  source     = "./notification"
  project_id = var.project_id
}
