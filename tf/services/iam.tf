locals {
  ml_bindings = {
    "roles/aiplatform.customCodeServiceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-aiplatform-cc.iam.gserviceaccount.com",
    ]
    "roles/aiplatform.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-aiplatform.iam.gserviceaccount.com",
      "serviceAccount:service-${var.project.number}@gcp-sa-aiplatform-ft.iam.gserviceaccount.com",
      "serviceAccount:service-${var.project.number}@gcp-sa-aiplatform-pie.iam.gserviceaccount.com",
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
  pubsub_identity = {
    member = "serviceAccount:service-${var.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
  }
  base_bindings = {
    "roles/bigquerydatatransfer.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com",
    ]
    "roles/compute.serviceAgent" = [
      data.google_compute_default_service_account.default.member,
    ]
    "roles/container.defaultNodeServiceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-gkenode.iam.gserviceaccount.com",
    ]
    "roles/container.nodeServiceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-gkenode.iam.gserviceaccount.com",
    ]
    "roles/container.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@container-engine-robot.iam.gserviceaccount.com",
    ]
    "roles/containeranalysis.ServiceAgent" = [
      "serviceAccount:service-${var.project.number}@container-analysis.iam.gserviceaccount.com",
    ]
    "roles/containerscanning.ServiceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-containerscanning.iam.gserviceaccount.com",
    ]
    "roles/eventarc.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-eventarc.iam.gserviceaccount.com",
    ]
    "roles/logging.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-logging.iam.gserviceaccount.com",
    ]
    "roles/memorystore.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-memorystore.iam.gserviceaccount.com",
    ]
    "roles/networkconnectivity.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-networkconnectivity.iam.gserviceaccount.com",
    ]
    "roles/networkmanagement.serviceAgent" = [
      "serviceAccount:service-${var.project.number}@gcp-sa-networkmanagement.iam.gserviceaccount.com",
    ]
    "roles/pubsub.serviceAgent" = [
      local.pubsub_identity.member,
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
    var.enable_ml ? local.ml_bindings : {},
    var.enable_redis ? local.redis_bindings : {},
  )
}
