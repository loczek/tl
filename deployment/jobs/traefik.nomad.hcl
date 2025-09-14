variable "image" {
  type    = string
  default = "traefik:3.5.0"
}

variable "cf_dns_api_token" {
  type = string
}

job "traefik" {
  type = "system"

  group "traefik" {
    network {
      port "traefik" {
        static = 8080
        to     = 8080
      }

      port "http" {
        static = 80
        to     = 80
      }

      port "https" {
        static = 443
        to     = 443
      }
    }

    task "traefik" {
      driver = "docker"

      config {
        image        = var.image
        network_mode = "host"
        ports = [
          "traefik",
          "http",
          "https"
        ]
        volumes = [
          "${NOMAD_TASK_DIR}/traefik.yaml:/etc/traefik/traefik.yaml",
        ]
      }

      env {
        CF_DNS_API_TOKEN = var.cf_dns_api_token
      }

      # Needed for the NOMAD_TOKEN env var
      identity {
        env         = true
        change_mode = "restart"
      }

      template {
        data        = file("./deployment/jobs/templates/traefik.yml.tmpl")
        destination = "${NOMAD_TASK_DIR}/traefik.yaml"
      }

      service {
        name     = "traefik-http"
        port     = "http"
        provider = "nomad"
        check {
          type     = "tcp"
          interval = "5s"
          timeout  = "1s"
        }
      }

      service {
        name     = "traefik-api"
        port     = "traefik"
        provider = "nomad"
        check {
          type     = "tcp"
          interval = "5s"
          timeout  = "1s"
        }
      }

      resources {
        cpu    = 200
        memory = 128
      }
    }
  }
}
