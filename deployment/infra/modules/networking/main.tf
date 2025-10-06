resource "aws_vpc" "main" {
  cidr_block                       = "10.0.0.0/16"
  assign_generated_ipv6_cidr_block = true

  tags = {
    "Name" = "tl-vpc"
  }
}

resource "aws_eip" "static_ip" {
  domain = "vpc"
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
