locals {
  location = "europe-west1"
}

resource "google_container_cluster" "main" {
  provider         = google-beta
  name             = var.name
  enable_autopilot = true

  min_master_version = var.container_engine_version.latest_master_version

  vertical_pod_autoscaling {
    enabled = true
  }

  release_channel {
    channel = var.release_channel_name
  }

  location   = local.location
  network    = var.network.name
  subnetwork = var.subnetwork.name
  ip_allocation_policy {
  }

  node_pool_auto_config {
    network_tags {
      tags = var.add_node_pool_network_tags
    }
  }

  private_cluster_config {
    enable_private_endpoint = false
    enable_private_nodes    = true
    master_ipv4_cidr_block  = var.master_cidr_block
    master_global_access_config {
      enabled = true
    }
  }

  gateway_api_config {
    channel = "CHANNEL_STANDARD"
  }

  deletion_protection = var.deletion_protection

  protect_config {
    workload_config {
      audit_mode = "BASIC"
    }
    workload_vulnerability_mode = "BASIC"
  }


  lifecycle {
    ignore_changes = [
      dns_config,
    ]
  }

  monitoring_config {
    enable_components = var.enable_monitoring_components
  }

}
