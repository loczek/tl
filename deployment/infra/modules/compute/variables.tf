variable "environment" {
  type = string
  default = "prod"
}

variable "ssh_public_key" {
  type = string
}

variable "public_eip_id" {
  type = string
}

variable "public_subnet" {
  type = list(string)
}

variable "private_subnet" {
  type = list(string)
}

variable "sg_public" {
  type = object({
    id = string
  })
}

variable "sg_internal" {
  type = object({
    id = string
  })
}

variable "sg_private" {
  type = object({
    id = string
  })
}
