# Full configuration options can be found at https://developer.hashicorp.com/nomad/docs/configuration
data_dir = "/opt/nomad/data"
bind_addr = "0.0.0.0"

acl {
  enabled = true
}

server {
  enabled          = true
  bootstrap_expect = 3

  server_join {
    retry_join = ["provider=aws tag_key=NomadServer tag_value=true"]
  }
}

client {
  enabled = true
}

telemetry {
  collection_interval        = "5s"
  publish_allocation_metrics = true
  publish_node_metrics       = true
  prometheus_metrics         = true
}
