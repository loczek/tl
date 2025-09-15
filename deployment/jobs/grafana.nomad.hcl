variable "domain" {
  type    = string
  default = "grafana.tiiinylink.com"
}

variable "image" {
  type    = string
  default = "grafana/grafana:latest"
}

job "grafana" {
  group "grafana" {
    network {
      port "http" {
        static = 9000
        to     = 3000
      }
    }

    service {
      name     = "grafana"
      port     = "http"
      provider = "nomad"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.grafana.tls=true",
        "traefik.http.routers.grafana.tls.certresolver=myresolver",
        "traefik.http.routers.grafana.entrypoints=websecure",
        "traefik.http.routers.grafana.rule=Host(`${var.domain}`)",
      ]
    }

    task "grafana" {
      driver = "docker"

      config {
        image = var.image
        ports = ["http"]
      }

      env {
        GF_SERVER_DOMAIN      = var.domain
        GF_SERVER_HTTP_PORT   = "${NOMAD_PORT_http}"
        GF_PATHS_PROVISIONING = "${NOMAD_TASK_DIR}/grafana/provisioning"
        GF_PATHS_CONFIG       = "${NOMAD_TASK_DIR}/grafana/grafana.ini"
      }

      template {
        data        = <<-EOF
        {{ with nomadVar "nomad/jobs/grafana" }}
        GF_SECURITY_ADMIN_USER = "{{ .username }}"
        GF_SECURITY_ADMIN_PASSWORD = "{{ .password }}"
        {{ end }}
        EOF
        env         = true
        destination = "${NOMAD_SECRETS_DIR}/pass.env"
      }

      template {
        data        = <<EOF

EOF
        destination = "${NOMAD_TASK_DIR}/grafana/grafana.ini"
      }

      template {
        data        = file("./deployment/jobs/templates/datasources.yml.tmpl")
        destination = "${NOMAD_TASK_DIR}/grafana/provisioning/datasources/datasources.yaml"
      }

      resources {
        cpu    = 200
        memory = 256
      }
    }
  }
}
