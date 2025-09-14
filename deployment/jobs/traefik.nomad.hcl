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
        image        = var.image
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
        destination = "local/traefik.yaml"
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

      resources {
        cpu    = 200
        memory = 128
      }
    }
  }
}
