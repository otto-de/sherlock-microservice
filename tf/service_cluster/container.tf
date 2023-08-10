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

  location   = "europe-west1"
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

  lifecycle {
    ignore_changes = [
      dns_config,
    ]
  }

}
