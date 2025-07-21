output "name" {
  description = "The name of the bucket."
  value       = google_storage_bucket.main.name
}

output "url" {
  description = "The URL of the bucket."
  value       = google_storage_bucket.main.url
}
