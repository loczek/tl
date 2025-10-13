job "otel-collector" {
  type = "system"

  group "otel-collector" {
    network {
      port "metrics" {
        static = 8888
        to     = 8888
      }
      port "prometheus" {
        static = 8889
        to     = 8889
      }
      port "otlp" {
        static = 4317
        to     = 4317
      }
      port "otlphttp" {
        static = 4318
        to     = 4318
      }
    }

    task "otel-collector" {
      driver = "docker"

      config {
        image = "otel/opentelemetry-collector-contrib:latest"
        entrypoint = [
          "/otelcol-contrib",
          "--config=${NOMAD_TASK_DIR}/otel/config.yaml",
        ]
        pid_mode   = "host"
        privileged = true

        ports = [
          "prometheus",
          "otlp",
          "otlphttp",
          "metrics"
        ]

        volumes = [
          "local/otel/config.yaml:/etc/otel/config.yaml",
          "/:/hostfs:ro,rslave",
        ]
      }

      # Needed for the NOMAD_TOKEN env var
      identity {
        env         = true
        change_mode = "restart"
      }

      env {
        HOST_DEV  = "/hostfs/dev"
        HOST_ETC  = "/hostfs/etc"
        HOST_PROC = "/hostfs/proc"
        HOST_RUN  = "/hostfs/run"
        HOST_SYS  = "/hostfs/sys"
        HOST_VAR  = "/hostfs/var"
      }

      template {
        data          = file("./deployment/jobs/templates/collector.yml.tmpl")
        change_mode   = "signal"
        change_signal = "SIGHUP"
        destination   = "${NOMAD_TASK_DIR}/otel/config.yaml"
      }

      resources {
        cpu    = 256
        memory = 128
      }

      service {
        name     = "otlp-collector"
        port     = "otlp"
        provider = "nomad"
        tags     = []
      }

      service {
        name     = "otlp-collector"
        port     = "otlphttp"
        provider = "nomad"
        tags     = ["otlphttp"]
      }

      service {
        name     = "otlp-collector"
        port     = "prometheus"
        provider = "nomad"
        tags     = ["prometheus"]
      }
    }
  }
}
