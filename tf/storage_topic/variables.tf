variable "bucket_configs" {
  type = map(object({
    object_name_prefix = string
  }))
}

variable "message_retention_duration" {
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
