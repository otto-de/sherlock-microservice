resource "google_pubsub_subscription" "default" {
  name = var.subscription_name

  topic                      = var.topic.id
  message_retention_duration = var.subscription_message_retention_duration
  enable_message_ordering    = var.enable_message_ordering
  ack_deadline_seconds       = var.ack_deadline_seconds
  filter                     = var.filter

  dead_letter_policy {
    dead_letter_topic = google_pubsub_topic.deadletter.id
  }
}

resource "google_pubsub_topic" "deadletter" {
  name = "${var.subscription_name}-deadletter"

  message_retention_duration = var.deadletter_message_retention_duration

  message_storage_policy {
    allowed_persistence_regions = [
      "europe-west1",
    ]
  }
}
