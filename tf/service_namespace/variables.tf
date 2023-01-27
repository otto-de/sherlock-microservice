variable "annotations" {
  type = map(string)
}

variable "labels" {
  type = map(string)
}

variable "name" {
  type        = string
  description = "Name of the Kubernetes namespace to create"
}

variable "pod_watchers" {
  type = list(object({
    metadata = list(object({
      name      = string
      namespace = string
    }))
  }))
  description = "Kubernetes Service Accounts that are allowed to watch for pods"
}
