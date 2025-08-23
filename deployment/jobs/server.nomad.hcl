job "server" {
  group "server" {
    count = 3

    network {
      port "http" {
        to = 3000
      }
    }

    service {
      name     = "server"
      port     = "http"
      provider = "nomad"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.server.rule=Host(`sho.rt`)",
        "traefik.http.routers.server.entrypoints=web",
      ]
    }

    task "server" {
      driver = "docker"

      config {
        image = "tl-server:local"
        ports = ["http"]
      }

      env {
        APP_ENV                     = "production"
        LOG_TO_STDOUT               = false
        OTEL_EXPORTER_OTLP_PROTOCOL = "http/protobuf"
        OTEL_METRIC_EXPORT_INTERVAL = "15000"
        OTEL_SERVICE_NAME           = "go-link-shortener"
      }

      template {
        data        = <<-EOF
        {{- range nomadService "postgres" -}}
        DATABASE_URL = "postgres://postgres:{{ with nomadVar "nomad/jobs/server" }}{{ .db_password }}{{ end }}@{{ .Address }}:{{ .Port }}/postgres"
        {{ end -}}

        {{- range nomadService "redis" -}}
        REDIS_URL = "redis://{{ .Address }}:{{ .Port }}/0"
        {{ end -}}

        {{- range nomadService "prometheus" -}}
        OTEL_EXPORTER_OTLP_METRICS_ENDPOINT = "http://{{ .Address }}:{{ .Port }}/api/v1/otlp/v1/metrics"
        {{ end -}}

        {{- range nomadService "loki" -}}
        OTEL_EXPORTER_OTLP_LOGS_ENDPOINT = "http://{{ .Address }}:{{ .Port }}/otlp/v1/logs"
        {{ end -}}

        {{- range nomadService "tempo-otlp-http" -}}
        OTEL_EXPORTER_OTLP_TRACES_ENDPOINT = "http://{{ .Address }}:{{ .Port }}/v1/traces"
        {{ end -}}
        EOF
        env         = true
        change_mode = "restart"
        destination = "local/other.env"
      }

      resources {
        cpu    = 100
        memory = 128
      }
    }
  }
}