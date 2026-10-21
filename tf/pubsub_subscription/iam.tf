data "google_iam_policy" "deadletter_topic" {

  binding {
    role = "roles/pubsub.publisher"
    members = [
      local.pubsub_identity.member,
    ]
  }
}

data "google_iam_policy" "subscription" {
  binding {
    role = "roles/pubsub.subscriber"
    members = distinct(concat(
      var.subscribers[*].member,
      # The service agent needs subscriber on the source subscription to forward dead letters.
      [local.pubsub_identity.member],
    ))
  }
  binding {
    role    = "roles/pubsub.viewer"
    members = var.viewers[*].member
  }
}

locals {
  pubsub_identity = {
    member = "serviceAccount:service-${data.google_project.default.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
  }
}

resource "google_pubsub_subscription_iam_policy" "default" {
  subscription = google_pubsub_subscription.default.id
  policy_data  = data.google_iam_policy.subscription.policy_data
}

resource "google_pubsub_topic_iam_policy" "deadletter" {
  topic       = google_pubsub_topic.deadletter.id
  policy_data = data.google_iam_policy.deadletter_topic.policy_data
}
