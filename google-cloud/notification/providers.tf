terraform {
  required_version = "= 1.14.3"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.17.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "= 2.7.1"
    }
  }
}
