
module "defs" {
  source = "../services"

  enable_artifactregistry = var.enable_artifactregistry
  enable_bigquery         = var.enable_bigquery
  enable_compute          = var.enable_compute
  enable_eventarc         = var.enable_eventarc
  enable_gke              = var.enable_gke
  enable_memorystore      = var.enable_memorystore
  enable_ml               = var.enable_ml
  enable_pubsub           = var.enable_pubsub
  enable_redis            = var.enable_redis
  project                 = var.project
}

# This just ensures that the correct service agent roles
# are assigned without touching any other roles
resource "google_project_iam_binding" "default" {
  for_each = module.defs.iam_bindings

  role    = each.key
  members = each.value
  project = var.project.name
}
