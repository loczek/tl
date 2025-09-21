variable "image" {
  type    = string
  default = "postgres:17.5-alpine"
}

job "postgres" {
  group "postgres" {
    network {
      port "postgres" {
        static = 5432
      }
    }

    service {
      name     = "postgres"
      port     = "postgres"
      provider = "nomad"

      check {
        name     = "alive"
        type     = "tcp"
        interval = "10s"
        timeout  = "2s"
      }
    }

    task postgres {
      driver = "docker"

      config {
        image = var.image
        ports = ["postgres"]
      }

      template {
        data        = <<-EOF
        POSTGRES_PASSWORD = "{{ with nomadVar "nomad/jobs/postgres" }}{{ .password }}{{ end }}"
        EOF
        env         = true
        destination = "${NOMAD_SECRETS_DIR}/password.env"
      }

      logs {
        max_files     = 5
        max_file_size = 15
      }

      resources {
        cpu    = 1000
        memory = 1024
      }
    }

    restart {
      attempts = 10
      interval = "5m"
      delay    = "25s"
      mode     = "delay"
    }
  }


  update {
    max_parallel     = 1
    min_healthy_time = "5s"
    healthy_deadline = "3m"
    auto_revert      = false
    canary           = 0
  }
}