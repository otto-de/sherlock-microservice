
resource "kubernetes_role_binding_v1" "pod_watcher" {
  metadata {
    name      = "pod-watchers"
    labels    = var.labels
    namespace = kubernetes_role_v1.pod_watcher.metadata.0.namespace
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "Role"
    name      = kubernetes_role_v1.pod_watcher.metadata.0.name
  }
  dynamic "subject" {
    for_each = var.pod_watchers
    content {
      kind = "ServiceAccount"

      name      = subject.value.metadata.0.name
      namespace = subject.value.metadata.0.namespace
    }
  }
}

resource "kubernetes_role_v1" "pod_watcher" {
  metadata {
    name      = "pod-watcher"
    namespace = kubernetes_namespace_v1.main.metadata.0.name
    labels    = var.labels
  }

  rule {
    api_groups = [""]
    resources  = ["pods"]
    verbs      = ["watch", "list", "get"]
  }
}
