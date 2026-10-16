
module "defs" {
  source = "../services"

  enable_ml    = var.enable_ml
  enable_redis = var.enable_redis
  project      = var.project
}

# This just ensures that the correct service agent roles
# are assigned without touching any other roles
resource "google_project_iam_binding" "default" {
  for_each = module.defs.iam_bindings

  role    = each.key
  members = each.value
  project = var.project.name
}
