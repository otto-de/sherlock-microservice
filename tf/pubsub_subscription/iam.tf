data "google_project" "default" {
}

data "google_iam_policy" "deadletter_topic" {

  binding {
    role    = "roles/pubsub.publisher"
    members = [
      "serviceAccount:service-${data.google_project.default.number}@gcp-sa-pubsub.iam.gserviceaccount.com",
    ]
  }
}

data "google_iam_policy" "subscription" {
  binding {
    role    = "roles/pubsub.subscriber"
    members = distinct(concat(
      var.subscribers[*].member,
      # The service agent needs subscriber on the source subscription to forward dead letters.
      ["serviceAccount:service-${data.google_project.default.number}@gcp-sa-pubsub.iam.gserviceaccount.com"],
    ))
  }
  binding {
    role = "roles/pubsub.viewer"
    members = var.viewers[*].member
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
