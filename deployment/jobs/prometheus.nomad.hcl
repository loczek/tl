job "prometheus" {
  group "prometheus" {

    network {
      port "http" {
        static = 9090
        to     = 9090
      }
    }

    service {
      name     = "prometheus"
      port     = "http"
      provider = "nomad"

      check {
        type     = "http"
        path     = "/-/healthy"
        interval = "3s"
        timeout  = "1s"
      }
    }


    task "prometheus" {
      driver = "docker"

      config {
        image = "prom/prometheus:v3.5.0"
        ports = ["http"]
        args = [
          "--web.enable-otlp-receiver",
          "--config.file=/etc/prometheus/config/prometheus.yml",
          "--storage.tsdb.path=/prometheus",
          "--web.listen-address=0.0.0.0:9090",
          # "--web.console.libraries=/usr/share/prometheus/console_libraries",
          # "--web.console.templates=/usr/share/prometheus/consoles"
        ]
        volumes = [
          "local/config:/etc/prometheus/config",
        ]
      }

      template {
        data        = file("./deployment/jobs/templates/prometheus.yml.tpl")
        destination = "local/config/prometheus.yml"
      }

      #       template {
      #         data = <<EOH
      # ---
      # global:
      #   scrape_interval: 30s
      #   evaluation_interval: 3s

      # scrape_configs:
      #   - job_name: prometheus
      #     static_configs:
      #     - targets:
      #       - 0.0.0.0:9090
      # EOH

      #         change_mode   = "signal"
      #         change_signal = "SIGHUP"
      #         destination   = "local/config/prometheus.yml"
      #       }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}