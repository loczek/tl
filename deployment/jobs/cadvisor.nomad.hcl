job "cadvisor" {
  type = "system"

  group "cadvisor" {
    task "cadvisor" {
      driver = "docker"

      config {
        image = "google/cadvisor:latest"
        args  = ["-docker_only"]
        volumes = [
          "/:/rootfs:ro",
          "/var/run:/var/run",
          "/sys:/sys:ro",
          "/var/lib/docker/:/var/lib/docker:ro",
          "/var/run/docker.sock:/var/run/docker.sock"
        ]
      }

      resources {
        cpu    = 100
        memory = 128
      }
    }
  }
}
