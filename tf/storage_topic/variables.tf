variable "bucket" {
  type = object({
    id = string
  })
}

variable "message_retention_duration" {
  type = string
}

variable "object_name_prefix" {
  type = string
}

variable "subscribers" {
  type = list(object({
    member = string
  }))
}

variable "topic_name" {
  type = string
}

variable "viewers" {
  type = list(object({
    member = string
  }))
}
