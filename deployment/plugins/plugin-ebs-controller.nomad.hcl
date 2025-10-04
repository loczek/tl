job "plugin-aws-ebs-controller" {
  group "controller" {
    constraint {
      attribute = "${meta.role}"
      operator  = "set_contains"
      value     = "ingress"
    }

    task "plugin" {
      driver = "docker"

      config {
        image        = "ecr-public.aws.com/ebs-csi-driver/aws-ebs-csi-driver:v1.50.0"
        network_mode = "host"
        args = [
          "controller",
          "--endpoint=unix://csi/csi.sock",
          "--logtostderr",
          "--v=5",
        ]
      }

      env {
        AWS_USE_DUALSTACK_ENDPOINT = "true"
      }

      csi_plugin {
        id        = "aws-ebs0"
        type      = "controller"
        mount_dir = "/csi"
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}
