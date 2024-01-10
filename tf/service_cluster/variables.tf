variable "add_node_pool_network_tags" {
  type    = list(string)
  default = []
}

variable "deletion_protection" {
  type        = bool
  description = <<EOF
Whether or not to allow Terraform to destroy the cluster.
Unless this field is set to false in Terraform state, a `terraform destroy` or `terraform apply` that would delete the cluster will fail.
EOF
}

variable "container_engine_version" {
  type = object({
    latest_master_version = string
  })
}

variable "enable_monitoring_components" {
  type        = list(string)
  default     = ["SYSTEM_COMPONENTS"]
  description = <<EOF
The GKE components exposing metrics.
Supported values include: `SYSTEM_COMPONENTS`, `APISERVER`, `SCHEDULER`, `CONTROLLER_MANAGER`, `STORAGE`, `HPA`, `POD`, `DAEMONSET`, `DEPLOYMENT` and `STATEFULSET`
EOF
}

variable "master_cidr_block" {
  type = string
}

variable "name" {
  type        = string
  description = "The name of the GKE cluster"
}

variable "network" {
  type = object({
    name = string
  })
}

variable "release_channel_name" {
  type    = string
  default = "STABLE"
}

variable "subnetwork" {
  type = object({
    name = string
  })
}

