moved {
  from = google_project_service.artifactregistry
  to   = google_project_service.required_apis["artifactregistry.googleapis.com"]
}

moved {
  from = google_project_service.billing_budgets
  to   = google_project_service.required_apis["billingbudgets.googleapis.com"]
}

moved {
  from = google_project_service.cloud_billing
  to   = google_project_service.required_apis["cloudbilling.googleapis.com"]
}

moved {
  from = google_project_service.cloudbuild
  to   = google_project_service.required_apis["cloudbuild.googleapis.com"]
}

moved {
  from = google_project_service.cloudfunctions
  to   = google_project_service.required_apis["cloudfunctions.googleapis.com"]
}

moved {
  from = google_project_service.run
  to   = google_project_service.required_apis["run.googleapis.com"]
}

moved {
  from = google_project_service.cloudresourcemanager
  to   = google_project_service.required_apis["cloudresourcemanager.googleapis.com"]
}

moved {
  from = google_project_service.eventarc
  to   = google_project_service.required_apis["eventarc.googleapis.com"]
}

moved {
  from = google_project_service.iamcredentials
  to   = google_project_service.required_apis["iamcredentials.googleapis.com"]
}

moved {
  from = google_project_service.iam
  to   = google_project_service.required_apis["iam.googleapis.com"]
}

moved {
  from = google_project_service.pubsub
  to   = google_project_service.required_apis["pubsub.googleapis.com"]
}

moved {
  from = google_project_service.secretmanager
  to   = google_project_service.required_apis["secretmanager.googleapis.com"]
}

moved {
  from = google_project_service.sts
  to   = google_project_service.required_apis["sts.googleapis.com"]
}
