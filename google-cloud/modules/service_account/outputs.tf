output "member" {
  description = "The email address of the created service account"
  value       = google_service_account.main.member
}

output "id" {
  description = "The ID of the created service account"
  value       = google_service_account.main.id
}

output "email" {
  description = "The email address of the created service account"
  value       = google_service_account.main.email
}
