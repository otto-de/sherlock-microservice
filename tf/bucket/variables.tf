
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

variable "bucket_viewers" {
  description = "IAM-style members who will be granted roles/storage.bucketViewer on bucket."
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
  type = object({
    retention_period = number
  })
  nullable    = true
  default     = null
  description = "Data retention policy for the bucket. retention_period is the minimum duration in seconds that objects must be retained. Null disables the policy."
}

variable "versioning_enabled" {
  description = "While set to true, versioning is fully enabled for this bucket."
  type        = bool
  default     = false
}

variable "object_ttl_days" {
  description = "Maximum age in days before an object in the bucket is deleted. Null disables the rule."
  type        = number
  nullable    = true
  default     = null
}
