
variable "name" {
  description = "The name of the bucket."
  type        = string
}

variable "location" {
  description = "The Google Cloud Storage location"
  type        = string
  default     = "europe-west1"
}

variable "force_destroy" {
  description = "When deleting a bucket, this boolean option will delete all contained objects. If you try to delete a bucket that contains objects, Terraform will fail that run."
  type        = bool
  default     = false
}


variable "admins" {
  description = "IAM-style members who will be granted roles/storage.objectAdmin on bucket."
  type        = list(string)
  default     = []
}

variable "viewers" {
  description = "IAM-style members who will be granted roles/storage.objectViewer on bucket."
  type        = list(string)
  default     = []
}

variable "users" {
  description = "IAM-style members who will be granted roles/storage.objectUser on bucket."
  type        = list(string)
  default     = []
}

variable "retention_policy" {
  type        = map(any)
  nullable    = true
  default     = null
  description = "Configuration of the bucket's data retention policy for how long objects in the bucket should be retained."
}
