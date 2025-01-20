data "google_iam_policy" "main_bucket" {
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

  force_destroy               = var.force_destroy
  public_access_prevention    = "enforced"
  uniform_bucket_level_access = true

  dynamic "retention_policy" {
    for_each = var.retention_policy == null ? [] : [var.retention_policy]
    content {
      retention_period = each.retention_period
    }
  }
}

resource "google_storage_bucket_iam_policy" "main" {
  bucket      = google_storage_bucket.main.name
  policy_data = data.google_iam_policy.main_bucket.policy_data
}
