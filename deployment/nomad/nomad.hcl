# Full configuration options can be found at https://developer.hashicorp.com/nomad/docs/configuration

data_dir = "/opt/nomad/data"

acl {
  enabled = true
}

server {
  enabled          = true
  bootstrap_expect = 1
}

client {
  enabled = true
  servers = ["127.0.0.1"]
}
