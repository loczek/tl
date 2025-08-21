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
        # "traefik.http.services.server.loadbalancer.server.port=3000"
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
        DATABASE_URL                = "postgres://postgres:example@postgres:5432/postgres"
        OTEL_EXPORTER_OTLP_PROTOCOL = "http/protobuf"
        OTEL_METRIC_EXPORT_INTERVAL = "15000"
        OTEL_SERVICE_NAME           = "go-link-shortener"
      }

      template {
        #         data        = <<EOH
        # REDIS_URL = "redis://{{ range nomadService "redis" }}{{ .Address }}:{{ .Port }}{{ end }}/0"
        #         EOH
        data        = <<EOH
        {{ with nomadService "redis" }}
          {{ with index . 0}}
            REDIS_URL = "redis://{{ .Address }}:{{ .Port }}/0"
          {{ end }}
        {{ end }}
        {{ with nomadService "prometheus" }}
          {{ with index . 0}}
            OTEL_EXPORTER_OTLP_METRICS_ENDPOINT = "http://{{ .Address }}:{{ .Port }}/api/v1/otlp/v1/metrics"
          {{ end }}
        {{ end }}
        {{ with nomadService "loki" }}
          {{ with index . 0}}
            OTEL_EXPORTER_OTLP_LOGS_ENDPOINT = "http://{{ .Address }}:{{ .Port }}/otlp/v1/logs"
          {{ end }}
        {{ end }}
        {{ with nomadService "tempo-otlp-http" }}
          {{ with index . 0}}
            OTEL_EXPORTER_OTLP_TRACES_ENDPOINT = "http://{{ .Address }}:{{ .Port }}/v1/traces"
          {{ end }}
        {{ end }}
        EOH
        env         = true
        change_mode = "restart"
        destination = "local/database.env"
      }

      resources {
        cpu    = 100
        memory = 128
      }
    }
  }
}