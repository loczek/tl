provider "aws" {
  region = "eu-central-1"
  default_tags {
    tags = {
      Environment = "staging"
    }
  }
}

module "networking" {
  source = "../modules/networking"
}

module "security" {
  source = "../modules/security"

  vpc_id = module.networking.vpc_id

  ip_whitelist = var.ip_whitelist
}

module "compute" {
  source = "../modules/compute"

  public_eip_id  = module.networking.public_eip_id
  public_subnet  = module.networking.public_subnet
  private_subnet = module.networking.private_subnet
  sg_public      = module.security.sg_public
  sg_internal    = module.security.sg_internal
  sg_private     = module.security.sg_private

  ssh_public_key = var.ssh_public_key
}
