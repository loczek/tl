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
        image = "ecr-public.aws.com/ebs-csi-driver/aws-ebs-csi-driver:v1.49.0"

        args = [
          "controller",
          "--endpoint=unix://csi/csi.sock",
          "--logtostderr",
          "--v=5",
        ]
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
