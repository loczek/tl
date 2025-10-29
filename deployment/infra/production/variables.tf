variable "ip_whitelist" {
  type        = string
  description = "IPv4 address that is able to connect to private services"
  sensitive   = true
}

variable "ssh_public_key" {
  type = string
}
