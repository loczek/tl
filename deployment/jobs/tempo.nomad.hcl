job "tempo" {
  group "tempo" {
    network {
      port "http" {
        to = 3200
      }

      port "grpc" {
        to = 9095
      }

      port "otlp_http" {
        to     = 4318
        static = 4318
      }

      port "otlp_grpc" {
        to     = 4317
        static = 4317
      }
    }


    service {
      name     = "tempo"
      port     = "http"
      provider = "nomad"

      check {
        type     = "http"
        path     = "/ready"
        interval = "5s"
        timeout  = "1s"
      }
    }

    service {
      name     = "tempo-otlp-http"
      port     = "otlp_http"
      provider = "nomad"
    }

    task "tempo" {
      driver = "docker"

      config {
        image = "grafana/tempo:2.8.1"
        ports = ["http", "grpc", "otlp_http", "otlp_grpc"]
        args = [
          "--config.file=/etc/tempo/config/tempo.yml",
        ]
        volumes = [
          "local/config:/etc/tempo/config",
        ]
      }

      resources {
        cpu    = 200
        memory = 256
      }

      template {
        data        = file("./deployment/jobs/templates/tempo.yml.tpl")
        destination = "local/config/tempo.yml"
      }

    }
  }
}