variable "artifact_registry_repository" {
  type = object({
    location      = string
    project       = string
    repository_id = string
  })
}

variable "build" {
  type = object({
    context    = string
    dockerfile = string
  })
}

variable "name" {
  type = string
}

variable "tag" {
  type = string
}
