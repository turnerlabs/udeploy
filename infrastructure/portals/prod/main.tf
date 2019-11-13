terraform {
  required_version = ">= 0.12.0"

  backend "s3" {
    region  = "us-east-1"
    profile = "{{PROFILE}}"
    bucket  = "{{tf-state-bucket-name}}"
    key     = "prod.udeploy.tfstate"
  }
}

variable "aws_profile" {
}

variable "region" {
  default = "us-east-1"
}

provider "aws" {
  version = ">= 1.46.0"
  region  = var.region
  profile = var.aws_profile
}

module "env" {
    source = "../../modules/portal"

    region = var.region
    aws_profile = var.aws_profile

    app = var.app
    environment = var.environment

    container_port = var.container_port
    health_check = var.health_check

    zone_name = var.zone_name
    record_name = var.domain

    vpc = var.vpc
    private_subnets = var.private_subnets
    public_subnets = var.public_subnets

    tags = var.tags

    saml_role = var.saml_role
    saml_users = var.saml_users

    # Optional. Uncomment the variable below if customizing the docker image is desired.
    # By default, a public image with the latest stable version of udeploy will be provided.
    # image = "{{CUSTOMIZED_IMAGE}}"

    config_path = var.config_path
}

output "alias_name" {
  value = module.env.dns_name
}

output "alias_zone_id" {
  value = module.env.zone_id
}

output "kms_key_id" {
  value = module.env.config_key_id
}

