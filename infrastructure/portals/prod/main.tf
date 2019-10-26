terraform {
  required_version = ">= 0.12.0"

  backend "s3" {
    region  = "us-east-1"
    profile = "{{PROFILE}}"
    bucket  = "tf-state-udeploy"
    key     = "prod.terraform.tfstate"
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

module "prod" {
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

    image = var.image

    config_path = var.config_path
}

