output "ip" {
  value = aws_eip.static_ip.public_ip
}

output "ssh" {
  value = "ssh ubuntu@${aws_eip.static_ip.public_ip}"
}

output "nomad_addr" {
  value = "export NOMAD_ADDR=http://${aws_eip.static_ip.public_ip}:4646"
}
}
