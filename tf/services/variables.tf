variable "enable_bigquery" {
  type    = bool
  default = false
}

variable "enable_eventarc" {
  type    = bool
  default = false
}

variable "enable_memorystore" {
  type    = bool
  default = false
}

variable "enable_ml" {
  type    = bool
  default = false
}

variable "enable_pubsub" {
  type    = bool
  default = false
}

variable "enable_redis" {
  type    = bool
  default = false
}

variable "project" {
  type = object({
    name   = string
    number = number
  })
}
