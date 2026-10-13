output "client_certificate" {
  value = google_container_cluster.main.master_auth.0.client_certificate
}

output "client_key" {
  value = google_container_cluster.main.master_auth.0.client_key
}

output "cluster_ca_certificate" {
  value = google_container_cluster.main.master_auth.0.cluster_ca_certificate
}

output "endpoint" {
  value = google_container_cluster.main.endpoint
}

output "location" {
  value = local.location
}

output "name" {
  value = var.name
}
