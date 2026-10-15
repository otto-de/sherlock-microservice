variable "enable_ml" {
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
