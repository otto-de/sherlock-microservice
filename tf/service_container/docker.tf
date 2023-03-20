locals {
  repository  = "${var.artifact_registry_repository.location}-docker.pkg.dev/${var.artifact_registry_repository.project}/${var.artifact_registry_repository.repository_id}"
  image       = "${local.repository}/${var.name}@${docker_registry_image.main.sha256_digest}"
  source_hash = sha1(join("", [for f in fileset(var.build.context, "**") : filesha1("${var.build.context}/${f}")]))
}

resource "docker_image" "main" {
  name = "${local.repository}/${var.name}:${var.tag}"

  build {
    context     = var.build.context
    dockerfile  = var.build.dockerfile
    pull_parent = true
  }

  triggers = {
    dir_sha1 = local.source_hash
  }
}

resource "docker_registry_image" "main" {
  name = docker_image.main.name

  triggers = {
    dir_sha1 = docker_image.main.image_id
  }

  lifecycle {
    create_before_destroy = true
  }
}
