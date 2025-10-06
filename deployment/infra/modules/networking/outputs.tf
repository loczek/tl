output "vpc_id" {
  value = aws_vpc.main.id
}

output "public_eip_id" {
  value = aws_eip.static_ip.id
}

output "public_ip" {
  value = aws_eip.static_ip.public_ip
}

output "public_subnet" {
  value = [aws_subnet.public-1.id, aws_subnet.public-2.id, aws_subnet.public-3.id]
}

output "private_subnet" {
  value = [aws_subnet.private-1.id, aws_subnet.private-2.id, aws_subnet.private-3.id]
}
