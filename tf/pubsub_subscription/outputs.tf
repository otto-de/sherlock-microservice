output "deadletter_topic" {
  value = google_pubsub_topic.deadletter
}

output "subscription" {
  value = google_pubsub_subscription.default
}
