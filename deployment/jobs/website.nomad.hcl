variable "domain" {
  type    = string
  default = "tiiinylink.com"
}

variable "image" {
  type    = string
  default = "loczek/tl-website"
}

job "website" {
  group "website" {
    constraint {
      attribute = "${meta.role}"
      operator  = "set_contains"
      value     = "ingress"
    }

    network {
      port "http" {
        to = 80
      }
    }

    service {
      name     = "website"
      port     = "http"
      provider = "nomad"

      tags = [
        "traefik.enable=true",
        "traefik.http.routers.website.rule=Host(`${var.domain}`)",
        "traefik.http.routers.website.entrypoints=web",
      ]
    }

    task "website" {
      driver = "docker"

      config {
        image = var.image
        ports = ["http"]
      }

      resources {
        cpu    = 200
        memory = 128
      }
    }
  }
}