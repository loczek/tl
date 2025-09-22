variable "image" {
  type    = string
  default = "redis:8.0-alpine"
}

job "redis" {
  group "redis" {
    constraint {
      attribute = "${meta.role}"
      operator  = "="
      value     = "ingress"
    }

    network {
      port "redis" {
        static = 6379
      }
    }

    # volume redis-data {
    #   type      = "host"
    #   source    = "redis"
    #   read_only = false
    # }

    service {
      name     = "redis"
      port     = "redis"
      provider = "nomad"

      # check {
      #   type     = "tcp"
      #   interval = "10s"
      #   timeout  = "2s"
      # }

      # check {
      #   name     = "redis_probe"
      #   type     = "tcp"
      #   interval = "10s"
      #   timeout  = "1s"
      # }
    }


    task "redis" {
      driver = "docker"

      config {
        image = var.image
        ports = ["redis"]
      }

      restart {
        # Nomad restarts a task on failure.  Attempts=0 → unlimited restarts.
        attempts = 0
        interval = "30s"
        delay    = "5s"
        mode     = "fail"
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}
