locals {
  services_config = yamldecode(file("${path.module}/gcp.yaml"))

  enabled = {
    artifactregistry = var.enable_artifactregistry
    bigquery         = var.enable_bigquery
    compute          = var.enable_compute || var.enable_gke
    eventarc         = var.enable_eventarc
    gke              = var.enable_gke
    memorystore      = var.enable_memorystore
    ml               = var.enable_ml
    pubsub           = var.enable_pubsub
    redis            = var.enable_redis
  }

  # service name => service config, for all services enabled by the feature flags
  active_service_configs = {
    for service in local.services_config.services : service.name => service
    if lookup(local.enabled, try(service.enabledBy, service.name), true)
  }

  # service name => google_project_service
  active_services = merge(google_project_service.this, google_project_service.dependent)

  # service name => (role => members)
  active_service_agent_roles = {
    for name, service in local.active_service_configs : name => {
      for agent in service.agents : agent.role => [
        for email in agent.emails : {
          member = "serviceAccount:${replace(email, "{project_number}", tostring(var.project.number))}"
        }
      ]
    }
    if can(service.agents)
  }

  # service name => [{member}]
  active_service_agents = {
    for name, roles in local.active_service_agent_roles : name => flatten(values(roles))
  }

  # role => members
  bindings = {
    for role, agents in merge(values(local.active_service_agent_roles)...) : role => agents[*].member
  }
}
