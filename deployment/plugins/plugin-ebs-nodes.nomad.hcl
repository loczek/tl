job "plugin-aws-ebs-nodes" {
  type = "system"

  group "nodes" {
    task "plugin" {
      driver = "docker"

      config {
        image = "ecr-public.aws.com/ebs-csi-driver/aws-ebs-csi-driver:v1.50.0"

        args = [
          "node",
          "--endpoint=unix://csi/csi.sock",
          "--logtostderr",
          "--v=5",
        ]

        privileged = true
      }

      csi_plugin {
        id             = "aws-ebs0"
        type           = "node"
        mount_dir      = "/csi"
        health_timeout = "60s"
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}
