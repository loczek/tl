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
        "traefik.http.routers.grafana.rule=Host(`grafana.short.com`)",
        "traefik.http.routers.grafana.entrypoints=web",
      ]
    }

    task "grafana" {
      driver = "docker"

      config {
        image = "grafana/grafana:latest"
        ports = ["http"]
      }

      resources {
        cpu    = 200
        memory = 256
      }

      env {
        # GF_LOG_LEVEL          = "DEBUG"
        # GF_LOG_MODE           = "console"
        GF_SERVER_HTTP_PORT   = "${NOMAD_PORT_http}"
        GF_PATHS_PROVISIONING = "/local/grafana/provisioning"
        GF_PATHS_CONFIG       = "/local/grafana/grafana.ini"
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
        destination = "/local/grafana/grafana.ini"
      }

      template {
        data        = file("./deployment/jobs/templates/datasources.yml.tpl")
        destination = "/local/grafana/provisioning/datasources/datasources.yaml"
      }
    }
  }
}