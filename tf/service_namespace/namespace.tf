resource "kubernetes_namespace_v1" "main" {
  metadata {
    annotations = var.annotations

    labels = var.labels

    name = var.name
  }
}
