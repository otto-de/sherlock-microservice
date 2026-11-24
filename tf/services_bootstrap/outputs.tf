output "dns" {
  value = module.defs.dns

  depends_on = [
    google_project_iam_binding.default,
  ]
}

output "logging" {
  value = module.defs.logging

  depends_on = [
    google_project_iam_binding.default,
  ]
}

output "secretmanager" {
  value = module.defs.secretmanager

  depends_on = [
    google_project_iam_binding.default,
  ]
}

output "servicenetworking" {
  value = module.defs.servicenetworking

  depends_on = [
    google_project_iam_binding.default,
  ]
}
