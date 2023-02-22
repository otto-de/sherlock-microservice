variable "add_node_pool_network_tags" {
  type    = list(string)
  default = []
}

variable "container_engine_version" {
  type = object({
    latest_master_version = string
  })
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

