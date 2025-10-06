provider "aws" {
  region = "eu-central-1"
}

module "networking" {
  source = "../modules/networking"
}

module "security" {
  source = "../modules/security"

  ip_whitelist = var.ip_whitelist
  vpc_id       = module.networking.vpc_id
}

resource "aws_iam_role" "nomad_server_role" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_policy" "nomad_server_policy" {
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect : "Allow",
        Action : [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ec2:DescribeAvailabilityZones",
          "ec2:DescribeInstances",
          "ec2:DescribeSnapshots",
          "ec2:DescribeTags",
          "ec2:DescribeVolumes",
          "ec2:DescribeVolumesModifications",
          "ec2:DescribeVolumeStatus"
        ],
        Resource : "*"
      },
      {
        Effect : "Allow",
        Action : [
          "ec2:CreateSnapshot",
          "ec2:ModifyVolume"
        ],
        Resource : "arn:aws:ec2:*:*:volume/*"
      },
      {
        Effect : "Allow",
        Action : [
          "ec2:AttachVolume",
          "ec2:DetachVolume"
        ],
        Resource : [
          "arn:aws:ec2:*:*:volume/*",
          "arn:aws:ec2:*:*:instance/*"
        ]
      },
      {
        Effect : "Allow",
        Action : [
          "ec2:CreateVolume",
          "ec2:EnableFastSnapshotRestores"
        ],
        Resource : "arn:aws:ec2:*:*:snapshot/*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "nomad_server_attach" {
  role       = aws_iam_role.nomad_server_role.name
  policy_arn = aws_iam_policy.nomad_server_policy.arn
}

resource "aws_iam_instance_profile" "nomad_server_profile" {
  name = "nomad-server-profile"
  role = aws_iam_role.nomad_server_role.name
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-arm64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_key_pair" "deployer" {
  key_name   = "tl-deployer-key"
  public_key = var.ssh_public_key
}

resource "aws_instance" "tl_monitoring" {
  private_ip        = "10.0.144.8"
  ami               = data.aws_ami.ubuntu.id
  availability_zone = "eu-central-1b"
  instance_type     = "t4g.small"
  subnet_id         = module.networking.private_subnet_2.id
  vpc_security_group_ids = [
    module.security.sg_private.id,
    module.security.sg_internal.id,
  ]
  key_name                    = aws_key_pair.deployer.key_name
  user_data_replace_on_change = true
  iam_instance_profile        = aws_iam_instance_profile.nomad_server_profile.name
  user_data = templatefile("${path.module}/../scripts/monitoring.sh.tmpl", {
    INSTALL_DOCKER_SCRIPT = file("${path.module}/../scripts/install-docker.sh")
    INSTALL_NOMAD_SCRIPT  = file("${path.module}/../scripts/install-nomad.sh")
    SETUP_NOMAD_SCRIPT    = file("${path.module}/../../nomad/monitoring.hcl")
  })

  root_block_device {
    encrypted = true
  }

  instance_market_options {
    market_type = "spot"

    spot_options {
      instance_interruption_behavior = "stop"
      spot_instance_type             = "persistent"
    }
  }

  metadata_options {
    http_put_response_hop_limit = 5
    http_protocol_ipv6          = "enabled"
    instance_metadata_tags      = "enabled"
  }

  tags = {
    "NomadServer" = true
  }
}

resource "aws_instance" "tl_database" {
  private_ip        = "10.0.144.16"
  ami               = data.aws_ami.ubuntu.id
  availability_zone = "eu-central-1b"
  instance_type     = "t4g.micro"
  subnet_id         = module.networking.private_subnet_2.id
  vpc_security_group_ids = [
    module.security.sg_private.id,
    module.security.sg_internal.id,
  ]
  key_name                    = aws_key_pair.deployer.key_name
  user_data_replace_on_change = true
  iam_instance_profile        = aws_iam_instance_profile.nomad_server_profile.name
  user_data = templatefile("${path.module}/../scripts/db.sh.tmpl", {
    INSTALL_DOCKER_SCRIPT = file("${path.module}/../scripts/install-docker.sh")
    INSTALL_NOMAD_SCRIPT  = file("${path.module}/../scripts/install-nomad.sh")
    SETUP_NOMAD_SCRIPT    = file("${path.module}/../../nomad/database.hcl")
  })

  root_block_device {
    encrypted = true
  }

  instance_market_options {
    market_type = "spot"

    spot_options {
      instance_interruption_behavior = "stop"
      spot_instance_type             = "persistent"
    }
  }

  metadata_options {
    instance_metadata_tags = "enabled"
  }

  tags = {
    "NomadServer" = true
  }
}

resource "aws_ebs_volume" "tl_database" {
  availability_zone = "eu-central-1b"
  size              = 20

  tags = {
    "Name" = "tl-ebs-database"
  }
}

resource "aws_instance" "tl_instance" {
  ami               = data.aws_ami.ubuntu.id
  availability_zone = "eu-central-1b"
  instance_type     = "t4g.micro" // "c6gd.medium"
  private_ip        = "10.0.16.8"
  subnet_id         = module.networking.public_subnet_2.id
  vpc_security_group_ids = [
    module.security.sg_public.id,
    module.security.sg_private.id,
    module.security.sg_internal.id,
  ]
  key_name                    = aws_key_pair.deployer.key_name
  user_data_replace_on_change = true
  iam_instance_profile        = aws_iam_instance_profile.nomad_server_profile.name
  user_data = templatefile("${path.module}/../scripts/ingress.sh.tmpl", {
    INSTALL_DOCKER_SCRIPT = file("${path.module}/../scripts/install-docker.sh")
    INSTALL_NOMAD_SCRIPT  = file("${path.module}/../scripts/install-nomad.sh")
    SETUP_NOMAD_SCRIPT    = file("${path.module}/../../nomad/ingress.hcl")
  })

  root_block_device {
    encrypted = true
  }

  instance_market_options {
    market_type = "spot"

    spot_options {
      instance_interruption_behavior = "stop"
      spot_instance_type             = "persistent"
    }
  }

  metadata_options {
    instance_metadata_tags = "enabled"
  }

  tags = {
    "NomadServer" = true
  }
}

resource "aws_eip" "static_ip" {
  domain = "vpc"
}

resource "aws_eip_association" "static_ip_assoc" {
  instance_id   = aws_instance.tl_instance.id
  allocation_id = aws_eip.static_ip.id
}
