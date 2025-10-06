
output "vpc_id" {
  value = aws_vpc.main.id
}

output "public_subnet_1" {
  value = aws_subnet.public-1
}

output "public_subnet_2" {
  value = aws_subnet.public-2
}

output "public_subnet_3" {
  value = aws_subnet.public-3
}

output "private_subnet_1" {
  value = aws_subnet.private-1
}

output "private_subnet_2" {
  value = aws_subnet.private-2
}

output "private_subnet_3" {
  value = aws_subnet.private-3
}
