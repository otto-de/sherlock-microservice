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
      tags = concat([
        "gke-keda-metrics",
      ], var.add_node_pool_network_tags)
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
}

resource "google_compute_firewall" "keda_allow_master" {
  name        = "gke-${var.name}-keda-metric-server"
  description = "KEDA Metric Server: A firewall rule to allow worker nodes communication with master"
  network     = var.network.name
  priority    = 1000
  direction   = "INGRESS"

  source_ranges = [
    var.master_cidr_block,
  ]
  target_tags = [
    "gke-keda-metrics",
  ]

  allow {
    protocol = "tcp"
    ports    = ["6443"]
  }

  log_config {
    metadata = "INCLUDE_ALL_METADATA"
  }
}
