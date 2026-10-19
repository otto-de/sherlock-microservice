resource "google_storage_notification" "default" {
  bucket             = var.bucket.id
  payload_format     = "JSON_API_V1"
  topic              = google_pubsub_topic.default.id
  object_name_prefix = var.object_name_prefix
  event_types        = ["OBJECT_FINALIZE"]

  depends_on = [
    google_pubsub_topic_iam_policy.default,
  ]
}
