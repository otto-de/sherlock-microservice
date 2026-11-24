output "dns" {
  value = {
    service = lookup(local.active_services, "dns", null)
  }
}

output "iam_bindings" {
  value = local.bindings

  depends_on = [
    google_project_service.this,
    google_project_service.dependent,
    terraform_data.service_identity,
  ]
}

output "logging" {
  value = {
    service = lookup(local.active_services, "logging", null)
  }
}

output "pubsub" {
  value = {
    agents = lookup(local.active_service_agents, "pubsub", [])
  }
}

output "secretmanager" {
  value = {
    service = lookup(local.active_services, "secretmanager", null)
  }
}

output "servicenetworking" {
  value = {
    agents  = lookup(local.active_service_agents, "servicenetworking", [])
    service = lookup(local.active_services, "servicenetworking", null)
  }
}
