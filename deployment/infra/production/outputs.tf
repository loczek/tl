output "ip" {
  value = module.networking.public_ip
}

output "ssh" {
  value = "ssh ubuntu@${module.networking.public_ip}"
}

output "nomad_addr" {
  value = "export NOMAD_ADDR=http://${module.networking.public_ip}:4646"
}

output "ssh_jump" {
  value = "ssh -o 'StrictHostKeyChecking no' -J ubuntu@${module.networking.public_ip} ubuntu@10.0.144.8"
}

output "db_volume_id" {
  value = module.compute.db_ebs_volume.id
}
