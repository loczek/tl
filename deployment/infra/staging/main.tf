provider "aws" {
  region = "eu-central-1"
}

resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"

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
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.128.0/20"
  availability_zone = "eu-central-1a"
  tags = {
    "Name" = "tl-private-1"
  }
}

resource "aws_subnet" "private-2" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.144.0/20"
  availability_zone = "eu-central-1b"
  tags = {
    "Name" = "tl-private-2"
  }
}

resource "aws_subnet" "private-3" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.160.0/20"
  availability_zone = "eu-central-1c"
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

resource "aws_security_group" "sg" {
  name        = "tl-security-group"
  description = "example"
  vpc_id      = aws_vpc.main.id
}

resource "aws_vpc_security_group_ingress_rule" "http" {
  security_group_id = aws_security_group.sg.id
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 80
  to_port           = 80
}

resource "aws_vpc_security_group_ingress_rule" "https" {
  security_group_id = aws_security_group.sg.id
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 433
  to_port           = 433
}

resource "aws_vpc_security_group_ingress_rule" "ssh" {
  security_group_id = aws_security_group.sg.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.my_ip
  from_port         = 22
  to_port           = 22
}

resource "aws_vpc_security_group_ingress_rule" "traefik-ui" {
  security_group_id = aws_security_group.sg.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.my_ip
  from_port         = 8080
  to_port           = 8080
}

resource "aws_vpc_security_group_ingress_rule" "nomad" {
  security_group_id = aws_security_group.sg.id
  ip_protocol       = "tcp"
  cidr_ipv4         = var.my_ip
  from_port         = 4646
  to_port           = 4646
}

resource "aws_vpc_security_group_ingress_rule" "nomad-server-tcp" {
  security_group_id = aws_security_group.sg.id
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 4648
  to_port           = 4648
}

resource "aws_vpc_security_group_ingress_rule" "nomad-server-udp" {
  security_group_id = aws_security_group.sg.id
  ip_protocol       = "udp"
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 4648
  to_port           = 4648
}

resource "aws_vpc_security_group_egress_rule" "all" {
  security_group_id = aws_security_group.sg.id
  ip_protocol       = -1
  cidr_ipv4         = "0.0.0.0/0"
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
  vpc_security_group_ids      = [aws_security_group.sg.id]
  key_name                    = aws_key_pair.deployer.key_name
  user_data_replace_on_change = true
  user_data = templatefile("${path.module}/../scripts/main.sh.tmpl", {
    GITHUB_TOKEN          = var.github_token
    GITHUB_USERNAME       = var.github_username
    INSTALL_DOCKER_SCRIPT = file("${path.module}/../scripts/install-docker.sh")
    INSTALL_NOMAD_SCRIPT  = file("${path.module}/../scripts/install-nomad.sh")
    SETUP_NOMAD_SCRIPT    = file("${path.module}/../../nomad/nomad.hcl")
  })
}
