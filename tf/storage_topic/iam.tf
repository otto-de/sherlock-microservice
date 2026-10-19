data "google_storage_project_service_account" "default" {
}

data "google_iam_policy" "topic" {
  binding {
    role = "roles/pubsub.publisher"
    members = [
      data.google_storage_project_service_account.default.member,
    ]
  }
  binding {
    role    = "roles/pubsub.subscriber"
    members = var.subscribers[*].member
  }
  binding {
    role    = "roles/pubsub.viewer"
    members = var.viewers[*].member
  }
}

resource "google_pubsub_topic_iam_policy" "default" {
  topic       = google_pubsub_topic.default.id
  policy_data = data.google_iam_policy.topic.policy_data
}
