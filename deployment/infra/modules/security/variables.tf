variable "environment" {
  type = string
  default = "prod"
}

variable "ip_whitelist" {
  type = string
}

variable "vpc_id" {
  type = string
}
