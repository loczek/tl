provider "aws" {
  region = "eu-central-1"
}

resource "aws_vpc" "main" {
  cidr_block                       = "10.0.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    "Name" = "tl-vpc"
  }
}

resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
  tags = {
    "Name" = "tl-internet-gateway"
  }
}

resource "aws_egress_only_internet_gateway" "egw" {
  vpc_id = aws_vpc.main.id
  tags = {
    "Name" = "tl-egress-only-internet-gateway"
  }
}

resource "aws_subnet" "public-1" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.0.0/20"
  availability_zone = "eu-central-1a"
  tags = {
    "Name" = "tl-public-1"
  }
}

resource "aws_subnet" "public-2" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.16.0/20"
  availability_zone = "eu-central-1b"
  tags = {
    "Name" = "tl-public-2"
  }
}

resource "aws_subnet" "public-3" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.32.0/20"
  availability_zone = "eu-central-1c"
  tags = {
    "Name" = "tl-public-3"
  }
}

resource "aws_subnet" "private-1" {
  vpc_id                          = aws_vpc.main.id
  cidr_block                      = "10.0.128.0/20"
  availability_zone               = "eu-central-1a"
  ipv6_cidr_block                 = cidrsubnet(aws_vpc.main.ipv6_cidr_block, 8, 0)
  assign_ipv6_address_on_creation = true

  tags = {
    "Name" = "tl-private-1"
  }
}

resource "aws_subnet" "private-2" {
  vpc_id                          = aws_vpc.main.id
  cidr_block                      = "10.0.144.0/20"
  availability_zone               = "eu-central-1b"
  ipv6_cidr_block                 = cidrsubnet(aws_vpc.main.ipv6_cidr_block, 8, 1)
  assign_ipv6_address_on_creation = true

  tags = {
    "Name" = "tl-private-2"
  }
}

resource "aws_subnet" "private-3" {
  vpc_id                          = aws_vpc.main.id
  cidr_block                      = "10.0.160.0/20"
  availability_zone               = "eu-central-1c"
  ipv6_cidr_block                 = cidrsubnet(aws_vpc.main.ipv6_cidr_block, 8, 2)
  assign_ipv6_address_on_creation = true

  tags = {
    "Name" = "tl-private-3"
  }
}

resource "aws_default_route_table" "private" {
  default_route_table_id = aws_vpc.main.default_route_table_id

  route {
    cidr_block = "10.0.0.0/16"
    gateway_id = "local"
  }

  route {
    ipv6_cidr_block = aws_vpc.main.ipv6_cidr_block
    gateway_id      = "local"
  }

  route {
    ipv6_cidr_block = "::/0"
    gateway_id      = aws_egress_only_internet_gateway.egw.id
  }

  lifecycle {
    ignore_changes = [route]
  }

  tags = {
    "Name" = "tl-private",
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "10.0.0.0/16"
    gateway_id = "local"
  }

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw.id
  }

  lifecycle {
    ignore_changes = [route]
  }

  tags = {
    "Name" = "tl-public",
  }
}

resource "aws_route_table_association" "a" {
  subnet_id      = aws_subnet.public-1.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "b" {
  subnet_id      = aws_subnet.public-2.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "c" {
  subnet_id      = aws_subnet.public-3.id
  route_table_id = aws_route_table.public.id
}

resource "aws_route_table_association" "d" {
  subnet_id      = aws_subnet.private-1.id
  route_table_id = aws_default_route_table.private.id
}

resource "aws_route_table_association" "e" {
  subnet_id      = aws_subnet.private-2.id
  route_table_id = aws_default_route_table.private.id
}

resource "aws_route_table_association" "f" {
  subnet_id      = aws_subnet.private-3.id
  route_table_id = aws_default_route_table.private.id
}

resource "aws_security_group" "public" {
  name        = "tl-sg-public"
  description = "HTTP/HTTPS ingress security group"
  vpc_id      = aws_vpc.main.id
}

resource "aws_vpc_security_group_ingress_rule" "http" {
  security_group_id = aws_security_group.public.id
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 80
  to_port           = 80
}

resource "aws_vpc_security_group_ingress_rule" "https" {
  security_group_id = aws_security_group.public.id
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 443
  to_port           = 443
}

resource "aws_vpc_security_group_egress_rule" "all" {
  security_group_id = aws_security_group.public.id
  ip_protocol       = -1
  cidr_ipv4         = "0.0.0.0/0"
}

resource "aws_security_group" "internal" {
  name        = "tl-sg-internal"
  description = "HTTP/HTTPS ingress security group"
  vpc_id      = aws_vpc.main.id
}

resource "aws_vpc_security_group_egress_rule" "internet" {
  security_group_id = aws_security_group.internal.id
  ip_protocol       = -1
  cidr_ipv6         = "::/0"
}

resource "aws_vpc_security_group_ingress_rule" "internal" {
  security_group_id = aws_security_group.internal.id
  ip_protocol       = -1
  cidr_ipv4         = "10.0.0.0/16"
}

resource "aws_vpc_security_group_egress_rule" "internal" {
  security_group_id = aws_security_group.internal.id
  ip_protocol       = -1
  cidr_ipv4         = "10.0.0.0/16"
}

resource "aws_security_group" "private" {
  name        = "tl-sg-private"
  description = "SSH security group"
  vpc_id      = aws_vpc.main.id
}

resource "aws_vpc_security_group_ingress_rule" "ssh" {
  security_group_id = aws_security_group.private.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.ip_whitelist
  from_port         = 22
  to_port           = 22
}

resource "aws_vpc_security_group_egress_rule" "ssh" {
  security_group_id = aws_security_group.private.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.ip_whitelist
  from_port         = 22
  to_port           = 22
}

resource "aws_vpc_security_group_ingress_rule" "traefik-ui" {
  security_group_id = aws_security_group.private.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.ip_whitelist
  from_port         = 8080
  to_port           = 8080
}

resource "aws_vpc_security_group_ingress_rule" "prometheus_ui" {
  security_group_id = aws_security_group.private.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.ip_whitelist
  from_port         = 9090
  to_port           = 9090
}

resource "aws_vpc_security_group_ingress_rule" "nomad" {
  security_group_id = aws_security_group.private.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.ip_whitelist
  from_port         = 4646
  to_port           = 4646
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
        Effect   = "Allow"
        Action   = ["ec2:DescribeInstances"]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = ["ec2:DescribeTags"]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = ["ec2:DescribeVolumes"]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = ["ec2:AttachVolume"]
        Resource = "*"
      },
      {
        Effect   = "Allow"
        Action   = ["ec2:DetachVolume"]
        Resource = "*"
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

resource "aws_instance" "tl_instance" {
  ami                         = data.aws_ami.ubuntu.id
  availability_zone           = "eu-central-1b"
  instance_type               = "c6gd.medium"
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.public-2.id
  vpc_security_group_ids = [
    aws_security_group.public.id,
    aws_security_group.private.id,
    aws_security_group.internal.id,
  ]
  key_name                    = aws_key_pair.deployer.key_name
  user_data_replace_on_change = true
  iam_instance_profile        = aws_iam_instance_profile.nomad_server_profile.name
  user_data = templatefile("${path.module}/../scripts/main.sh.tmpl", {
    GITHUB_TOKEN          = var.github_token
    GITHUB_USERNAME       = var.github_username
    INSTALL_DOCKER_SCRIPT = file("${path.module}/../scripts/install-docker.sh")
    INSTALL_NOMAD_SCRIPT  = file("${path.module}/../scripts/install-nomad.sh")
    SETUP_NOMAD_SCRIPT    = file("${path.module}/../../nomad/nomad.hcl")
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

