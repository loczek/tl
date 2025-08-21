job "website" {
  group "website" {
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
        "traefik.http.routers.website.rule=Host(`short.com`)",
        "traefik.http.routers.website.entrypoints=web",
      ]
    }

    task "website" {
      driver = "docker"

      config {
        image = "tl-website:local"
        ports = ["http"]
      }
    }
  }
}