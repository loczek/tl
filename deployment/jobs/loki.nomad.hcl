variable "image" {
  type    = string
  default = "grafana/loki:latest"
}

job "loki" {
  group "loki" {
    constraint {
      attribute = "${meta.role}"
      operator  = "="
      value     = "monitoring"
    }

    network {
      port "http" {
        static = 3100
        to     = 3100
      }
    }

    service {
      name     = "loki"
      port     = "http"
      provider = "nomad"
    }

    task "loki" {
      driver = "docker"

      config {
        image = var.image
        ports = ["http"]
        volumes = [
          # "loki-data:/loki",
          "local/config/loki.yaml:/etc/loki/loki.yaml",
        ]
      }

      template {
        data        = file("./loki/loki.yaml")
        destination = "${NOMAD_TASK_DIR}/config/loki.yaml"
      }

      resources {
        cpu    = 200
        memory = 256
      }
    }
  }
}
