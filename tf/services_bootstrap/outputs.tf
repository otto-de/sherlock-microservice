output "dns_service" {
  value = google_project_service.dns

  depends_on = [
    google_project_iam_binding.default,
  ]
}

output "logging_service" {
  value = google_project_service.logging

  depends_on = [
    google_project_iam_binding.default,
  ]
}

output "secretmanager_service" {
  value = google_project_service.secretmanager

  depends_on = [
    google_project_iam_binding.default,
  ]
}

output "servicenetworking_service" {
  value = google_project_service.servicenetworking

  depends_on = [
    google_project_iam_binding.default,
  ]
}
