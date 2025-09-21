variable "ip_whitelist" {
  type        = string
  description = "IPv4 address that is able to connect to private services"
  sensitive   = true
}

variable "github_token" {
  type      = string
  sensitive = true
}

variable "ssh_public_key" {
  type = string
}

variable "db_username" {
  type = string
}

variable "db_password" {
  type      = string
  sensitive = true
}
