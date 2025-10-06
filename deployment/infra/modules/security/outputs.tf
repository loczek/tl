output "sg_public" {
  value = aws_security_group.public
}

output "sg_private" {
  value = aws_security_group.private
}

output "sg_internal" {
  value = aws_security_group.internal
}
