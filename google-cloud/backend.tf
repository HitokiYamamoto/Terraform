terraform {
  required_version = "1.14.4"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.16.0"
    }
  }

  backend "gcs" {
    bucket = "terraform-backend-bucket-bue8pj"
    prefix = "state/google-cloud"
  }
}

provider "google" {
  project = var.project_id
  default_labels = {
    "env" = "development"
  }
}
