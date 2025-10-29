namespace = "default"
name      = "redis-volume"
type      = "host"

plugin_id = "mkdir"

constraint {
  attribute = "${meta.role}"
  operator  = "set_contains"
  value     = "database"
}

capability {
  access_mode     = "single-node-single-writer"
  attachment_mode = "file-system"
}
