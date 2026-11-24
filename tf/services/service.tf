# Services are configured in gcp.yaml.
resource "google_project_service" "this" {
  for_each = { for name, service in local.active_service_configs : name => service if !can(service.dependsOn) }

  service = "${each.key}.googleapis.com"

  disable_on_destroy = false
}

# Services that may only be enabled after the services they depend on.
resource "google_project_service" "dependent" {
  for_each = { for name, service in local.active_service_configs : name => service if can(service.dependsOn) }

  service = "${each.key}.googleapis.com"

  disable_on_destroy = false

  depends_on = [
    google_project_service.this,
  ]
}

resource "terraform_data" "service_identity" {
  for_each = { for name, service in local.active_service_configs : name => service if try(service.createIdentity, false) }

  provisioner "local-exec" {
    command = "gcloud beta services identity create --service=${local.active_services[each.key].service} --project=${var.project.name}"
  }

  triggers_replace = [var.project.name]

  depends_on = [
    google_project_service.this,
    google_project_service.dependent,
  ]
}
