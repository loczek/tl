namespace = "default"
name      = "grafana-volume"
type      = "host"

plugin_id = "mkdir"

constraint {
  attribute = "${meta.role}"
  operator  = "set_contains"
  value     = "monitoring"
}

capability {
  access_mode     = "single-node-single-writer"
  attachment_mode = "file-system"
}
