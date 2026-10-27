data "google_iam_policy" "main_bucket" {
  dynamic "binding" {
    for_each = var.bucket_viewers == [] ? [] : [var.bucket_viewers]
    content {
      role    = "roles/storage.bucketViewer"
      members = binding.value
    }
  }

  dynamic "binding" {
    for_each = var.admins == [] ? [] : [var.admins]
    content {
      role    = "roles/storage.objectAdmin"
      members = binding.value
    }
  }

  dynamic "binding" {
    for_each = var.users == [] ? [] : [var.users]
    content {
      role    = "roles/storage.objectUser"
      members = binding.value
    }
  }


  dynamic "binding" {
    for_each = var.viewers == [] ? [] : [var.viewers]
    content {
      role    = "roles/storage.objectViewer"
      members = binding.value
    }
  }
}

resource "google_storage_bucket" "main" {
  name = var.name

  location = var.location
  project  = var.project

  force_destroy               = var.force_destroy
  public_access_prevention    = "enforced"
  uniform_bucket_level_access = true
  storage_class               = var.storage_class
  labels                      = var.labels

  versioning {
    enabled = var.versioning_enabled
  }

  dynamic "soft_delete_policy" {
    for_each = var.soft_delete_retention_seconds == null ? [] : [var.soft_delete_retention_seconds]
    content {
      retention_duration_seconds = soft_delete_policy.value
    }
  }

  dynamic "retention_policy" {
    for_each = var.retention_policy == null ? [] : [var.retention_policy]
    content {
      retention_period = retention_policy.value.retention_period
    }
  }

  dynamic "lifecycle_rule" {
    for_each = var.object_ttl_days == null ? [] : [var.object_ttl_days]
    content {
      action {
        type = "Delete"
      }
      condition {
        age = lifecycle_rule.value
      }
    }
  }

  lifecycle {
    prevent_destroy = var.prevent_destroy
  }
}

resource "google_storage_bucket_iam_policy" "main" {
  bucket      = google_storage_bucket.main.name
  policy_data = data.google_iam_policy.main_bucket.policy_data
}
