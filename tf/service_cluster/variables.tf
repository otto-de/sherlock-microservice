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

variable "subnetwork" {
  type = object({
    name = string
  })
}

