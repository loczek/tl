resource "aws_security_group" "public" {
  name        = "tl-sg-public"
  description = "HTTP/HTTPS ingress security group"
  vpc_id      = var.vpc_id
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
  vpc_id      = var.vpc_id
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
  vpc_id      = var.vpc_id
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
