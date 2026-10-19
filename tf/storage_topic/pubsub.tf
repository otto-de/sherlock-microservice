resource "google_pubsub_topic" "default" {
  name = var.topic_name

  message_retention_duration = var.message_retention_duration

  message_storage_policy {
    allowed_persistence_regions = [
      "europe-west1",
    ]
  }
}
