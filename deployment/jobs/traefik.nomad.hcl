job "traefik" {
  type = "system"

  group "traefik" {
    network {
      port "api" {
        static = 1936
        to     = 1936
      }

      port "http" {
        static = 8080
        to     = 8080
      }

      port "https" {
        static = 8443
        to     = 8443
      }
    }

    task "traefik" {
      driver = "docker"

      config {
        image        = "traefik:3.5.0"
        network_mode = "host"
        ports = [
          "api",
          "http",
          "https"
        ]
        volumes = [
          "local/traefik.yaml:/etc/traefik/traefik.yaml",
        ]
      }

      template {
        data        = file("./deployment/jobs/templates/traefik.yml.tpl")
        destination = "local/traefik.yaml"
      }

      resources {
        cpu    = 200
        memory = 256
      }

      service {
        name     = "traefik-http"
        port     = "http"
        provider = "nomad"
        check {
          type     = "tcp"
          interval = "3s"
          timeout  = "1s"
        }
      }

      service {
        name     = "traefik-api"
        port     = "api"
        provider = "nomad"
        check {
          type     = "tcp"
          interval = "3s"
          timeout  = "1s"
        }
      }
    }
  }
}