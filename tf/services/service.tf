data "google_compute_default_service_account" "default" {
  depends_on = [
    google_project_service.compute
  ]
}

resource "google_project_service" "aiplatform" {
  service = "aiplatform.googleapis.com"

  disable_on_destroy = false

  lifecycle {
    enabled = var.enable_ml
  }
}

resource "google_project_service" "bigquerydatatransfer" {
  service = "bigquerydatatransfer.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "cloudbuild" {
  service = "cloudbuild.googleapis.com"

  disable_on_destroy = false

  lifecycle {
    enabled = var.enable_ml
  }
}

resource "google_project_service" "clouderrorreporting" {
  service = "clouderrorreporting.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "cloudresourcemanager" {
  service = "cloudresourcemanager.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "cloudtrace" {
  service = "cloudtrace.googleapis.com"

  disable_on_destroy = false
}

/*
 * Terraform uses Cloud Resource Manager for everything.
 * Thus we do not explicitly express the dependency.
resource "google_project_service" "cloudresourcemanager" {
  service = "cloudresourcemanager.googleapis.com"

  disable_on_destroy = false
}
*/

resource "google_project_service" "compute" {
  service = "compute.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "container" {
  service = "container.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "containeranalysis" {
  service = "containeranalysis.googleapis.com"

  disable_on_destroy = false

  depends_on = [
    google_project_service.pubsub,
  ]
}

resource "google_project_service" "containerscanning" {
  service = "containerscanning.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "containersecurity" {
  service = "containersecurity.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "dns" {
  service = "dns.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "eventarc" {
  service = "eventarc.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "eventarcpublishing" {
  service = "eventarcpublishing.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "iam" {
  service = "iam.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "iamcredentials" {
  service = "iamcredentials.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "iap" {
  service = "iap.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "logging" {
  service = "logging.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "memorystore" {
  service = "memorystore.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "ml" {
  service = "ml.googleapis.com"

  disable_on_destroy = false

  lifecycle {
    enabled = var.enable_ml
  }
}

resource "google_project_service" "networkconnectivity" {
  service = "networkconnectivity.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "networkmanagement" {
  service = "networkmanagement.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "notebooks" {
  service = "notebooks.googleapis.com"

  disable_on_destroy = false

  lifecycle {
    enabled = var.enable_ml
  }
}

resource "google_project_service" "observability" {
  service = "observability.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "pubsub" {
  service = "pubsub.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "redis" {
  service = "redis.googleapis.com"

  disable_on_destroy = false

  depends_on = [
    google_project_service.compute,
  ]

  lifecycle {
    enabled = var.enable_redis
  }
}

resource "google_project_service" "secretmanager" {
  service = "secretmanager.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "servicenetworking" {
  service = "servicenetworking.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "serviceusage" {
  service = "serviceusage.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "vpcaccess" {
  service = "vpcaccess.googleapis.com"

  disable_on_destroy = false
}
