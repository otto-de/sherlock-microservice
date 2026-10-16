output "dns_service" {
  value = module.defs.dns_service

  depends_on = [
    google_project_iam_binding.default,
  ]
}

output "logging_service" {
  value = module.defs.logging_service

  depends_on = [
    google_project_iam_binding.default,
  ]
}

output "secretmanager_service" {
  value = module.defs.secretmanager_service

  depends_on = [
    google_project_iam_binding.default,
  ]
}

output "servicenetworking_service" {
  value = module.defs.servicenetworking_service

  depends_on = [
    google_project_iam_binding.default,
  ]
}
