locals {
  artifactregistry_bindings = {
    "roles/artifactregistry.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-artifactregistry.iam.gserviceaccount.com"
    ]
    "roles/containeranalysis.ServiceAgent" = [
      "serviceAccount:service-${var.project.number}@container-analysis.iam.gserviceaccount.com",
    ]
    "roles/containerscanning.ServiceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-containerscanning.iam.gserviceaccount.com",
    ]
  }
  bigquery_bindings = {
    "roles/bigquerydatatransfer.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com",
    ]
  }
  compute_bindings = {
    "roles/compute.serviceAgent" = [
      data.google_compute_default_service_account.default.member,
    ]
  }
  eventarc_bindings = {
    "roles/eventarc.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-eventarc.iam.gserviceaccount.com",
    ]
  }
  gke_bindings = {
    "roles/container.defaultNodeServiceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-gkenode.iam.gserviceaccount.com",
    ]
    "roles/container.nodeServiceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-gkenode.iam.gserviceaccount.com",
    ]
    "roles/container.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@container-engine-robot.iam.gserviceaccount.com",
    ]
  }
  memorystore_bindings = {
    "roles/memorystore.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-memorystore.iam.gserviceaccount.com",
    ]
  }
  ml_bindings = {
    "roles/aiplatform.customCodeServiceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-aiplatform-cc.iam.gserviceaccount.com",
    ]
    "roles/aiplatform.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-aiplatform.iam.gserviceaccount.com",
      "serviceAccount:service-${var.project.number}@gcp-sa-aiplatform-ft.iam.gserviceaccount.com",
      "serviceAccount:service-${var.project.number}@gcp-sa-aiplatform-pie.iam.gserviceaccount.com",
    ]
    "roles/aiplatform.telemetryServiceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-vertex-telemetry.iam.gserviceaccount.com",
    ]
    "roles/cloudbuild.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-cloudbuild.iam.gserviceaccount.com",
    ]
    "roles/ml.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@cloud-ml.google.com.iam.gserviceaccount.com",
    ]
    "roles/notebooks.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-notebooks.iam.gserviceaccount.com",
    ]
  }
  redis_bindings = {
    "roles/redis.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@cloud-redis.iam.gserviceaccount.com",
    ]
  }
  pubsub_bindings = {
    "roles/pubsub.serviceAgent" = [
      local.pubsub_identity.member,
    ]
  }
  pubsub_identity = {
    member = "serviceAccount:service-${var.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
  }
  base_bindings = {
    "roles/logging.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-logging.iam.gserviceaccount.com",
    ]

    "roles/networkconnectivity.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-networkconnectivity.iam.gserviceaccount.com",
    ]
    "roles/networkmanagement.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-networkmanagement.iam.gserviceaccount.com",
    ]
    "roles/servicenetworking.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@service-networking.iam.gserviceaccount.com",
    ]
    "roles/vpcaccess.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-vpcaccess.iam.gserviceaccount.com",
    ]
  }
  bindings = merge(
    local.base_bindings,
    var.enable_artifactregistry ? local.artifactregistry_bindings : {},
    var.enable_bigquery ? local.bigquery_bindings : {},
    var.enable_compute ? local.compute_bindings : {},
    var.enable_eventarc ? local.eventarc_bindings : {},
    var.enable_gke ? local.gke_bindings : {},
    var.enable_memorystore ? local.memorystore_bindings : {},
    var.enable_ml ? local.ml_bindings : {},
    var.enable_pubsub ? local.pubsub_bindings : {},
    var.enable_redis ? local.redis_bindings : {},
  )
}
