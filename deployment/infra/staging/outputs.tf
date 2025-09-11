output "ip" {
  value = aws_instance.tl_instance.public_ip
}

output "ssh" {
  value = "ssh ubuntu@${aws_instance.tl_instance.public_ip}"
}

output "nomad_addr" {
  value = "export NOMAD_ADDR=http://${aws_instance.tl_instance.public_ip}:4646"
}
