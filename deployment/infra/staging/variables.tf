variable "my_ip" {
  type      = string
  sensitive = true
}

variable "github_token" {
  type      = string
  sensitive = true
}

variable "ssh_public_key" {
  type = string
}
