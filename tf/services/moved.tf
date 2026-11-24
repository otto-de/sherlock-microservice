moved {
  from = google_project_service.aiplatform
  to   = google_project_service.this["aiplatform"]
}

moved {
  from = google_project_service.artifactregistry
  to   = google_project_service.this["artifactregistry"]
}

moved {
  from = google_project_service.bigquerydatatransfer
  to   = google_project_service.this["bigquerydatatransfer"]
}

moved {
  from = google_project_service.cloudbuild
  to   = google_project_service.this["cloudbuild"]
}

moved {
  from = google_project_service.clouderrorreporting
  to   = google_project_service.this["clouderrorreporting"]
}

moved {
  from = google_project_service.cloudresourcemanager
  to   = google_project_service.this["cloudresourcemanager"]
}

moved {
  from = google_project_service.cloudtrace
  to   = google_project_service.this["cloudtrace"]
}

moved {
  from = google_project_service.compute
  to   = google_project_service.this["compute"]
}

moved {
  from = google_project_service.container
  to   = google_project_service.dependent["container"]
}

moved {
  from = google_project_service.containeranalysis
  to   = google_project_service.this["containeranalysis"]
}

moved {
  from = google_project_service.containerscanning
  to   = google_project_service.this["containerscanning"]
}

moved {
  from = google_project_service.containersecurity
  to   = google_project_service.this["containersecurity"]
}

moved {
  from = google_project_service.dns
  to   = google_project_service.this["dns"]
}

moved {
  from = google_project_service.eventarc
  to   = google_project_service.this["eventarc"]
}

moved {
  from = google_project_service.eventarcpublishing
  to   = google_project_service.this["eventarcpublishing"]
}

moved {
  from = google_project_service.iam
  to   = google_project_service.this["iam"]
}

moved {
  from = google_project_service.iamcredentials
  to   = google_project_service.this["iamcredentials"]
}

moved {
  from = google_project_service.iap
  to   = google_project_service.this["iap"]
}

moved {
  from = google_project_service.logging
  to   = google_project_service.this["logging"]
}

moved {
  from = google_project_service.memorystore
  to   = google_project_service.this["memorystore"]
}

moved {
  from = google_project_service.ml
  to   = google_project_service.this["ml"]
}

moved {
  from = google_project_service.networkconnectivity
  to   = google_project_service.this["networkconnectivity"]
}

moved {
  from = google_project_service.networkmanagement
  to   = google_project_service.this["networkmanagement"]
}

moved {
  from = google_project_service.notebooks
  to   = google_project_service.this["notebooks"]
}

moved {
  from = google_project_service.observability
  to   = google_project_service.this["observability"]
}

moved {
  from = google_project_service.pubsub
  to   = google_project_service.this["pubsub"]
}

moved {
  from = google_project_service.redis
  to   = google_project_service.this["redis"]
}

moved {
  from = google_project_service.secretmanager
  to   = google_project_service.this["secretmanager"]
}

moved {
  from = google_project_service.servicenetworking
  to   = google_project_service.this["servicenetworking"]
}

moved {
  from = google_project_service.serviceusage
  to   = google_project_service.this["serviceusage"]
}

moved {
  from = google_project_service.vpcaccess
  to   = google_project_service.this["vpcaccess"]
}

moved {
  from = google_project_service.workloadidentity
  to   = google_project_service.this["workloadidentity"]
}
