terraform {
  required_version = "1.14.4"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.18.0"
    }
  }
}


provider "google" {
  project = var.project_id
  default_labels = {
    "env" = "development"
  }
}
