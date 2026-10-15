locals {
  ml_bindings = [{
    role = "roles/aiplatform.customCodeServiceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-aiplatform-cc.iam.gserviceaccount.com",
    ] }, {
    role = "roles/aiplatform.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-aiplatform-ft.iam.gserviceaccount.com",
    ] }, {
    role    = "roles/cloudbuild.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-cloudbuild.iam.gserviceaccount.com",
    ] }, {
    role = "roles/notebooks.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-notebooks.iam.gserviceaccount.com",
    ] }, {
    role = "roles/ml.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@cloud-ml.google.com.iam.gserviceaccount.com",
    ]
  }]
  redis_bindings = [{
    role = "roles/redis.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@cloud-redis.iam.gserviceaccount.com",
    ]
  }]
  bindings = concat([{
    role = "roles/bigquerydatatransfer.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-bigquerydatatransfer.iam.gserviceaccount.com",
    ] }, {
    role = "roles/container.defaultNodeServiceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-gkenode.iam.gserviceaccount.com",
    ] }, {
    role = "roles/container.nodeServiceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-gkenode.iam.gserviceaccount.com",
    ] }, {
    role = "roles/container.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@container-engine-robot.iam.gserviceaccount.com",
    ] }, {
    role = "roles/containeranalysis.ServiceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@container-analysis.iam.gserviceaccount.com",
    ] }, {
    role = "roles/containerscanning.ServiceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-containerscanning.iam.gserviceaccount.com",
    ] }, {
    role = "roles/compute.serviceAgent"
    members = [
      data.google_compute_default_service_account.default.member,
    ] }, {
    role = "roles/eventarc.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-eventarc.iam.gserviceaccount.com",
    ] }, {
    role = "roles/logging.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-logging.iam.gserviceaccount.com",
    ] }, {
    role = "roles/networkmanagement.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-networkmanagement.iam.gserviceaccount.com",
    ] }, {
    role = "roles/networkconnectivity.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-networkconnectivity.iam.gserviceaccount.com",
    ] }, {
    role = "roles/memorystore.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-memorystore.iam.gserviceaccount.com",
    ] }, {
    role = "roles/pubsub.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com",
    ] }, {
    role = "roles/servicenetworking.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@service-networking.iam.gserviceaccount.com",
    ] }, {
    role = "roles/vpcaccess.serviceAgent"
    members = [
      "serviceAccount:service-${var.project.number}@gcp-sa-vpcaccess.iam.gserviceaccount.com",
    ]
    }],
    var.enable_ml ? local.ml_bindings : [],
    var.enable_redis ? local.redis_bindings : [],
  )
}

# This just ensures that the correct service agent roles
# are assigned without touching any other roles
resource "google_project_iam_binding" "default" {
  for_each = { for b in local.bindings : b.role => b }

  role    = each.value.role
  members = each.value.members
  project = var.project.name

  depends_on = [
    google_project_service.aiplatform,
    google_project_service.bigquerydatatransfer,
    google_project_service.cloudbuild,
    google_project_service.clouderrorreporting,
    google_project_service.cloudresourcemanager,
    google_project_service.cloudtrace,
    google_project_service.compute,
    google_project_service.container,
    google_project_service.containeranalysis,
    google_project_service.containerscanning,
    google_project_service.containersecurity,
    google_project_service.dns,
    google_project_service.eventarc,
    google_project_service.eventarcpublishing,
    google_project_service.iam,
    google_project_service.iamcredentials,
    google_project_service.iap,
    google_project_service.logging,
    google_project_service.memorystore,
    google_project_service.ml,
    google_project_service.networkconnectivity,
    google_project_service.networkmanagement,
    google_project_service.notebooks,
    google_project_service.observability,
    google_project_service.pubsub,
    google_project_service.redis,
    google_project_service.secretmanager,
    google_project_service.servicenetworking,
    google_project_service.serviceusage,
    google_project_service.vpcaccess,
  ]
}
