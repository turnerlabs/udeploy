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
    source = "github.com/turnerlabs/udeploy//infrastructure/modules/portal?ref=v0.32.2-rc"

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

    config_path = var.config_path

    parameter_store = var.parameter_store
}

# KMS Key used to encrypt the portal 
# SSM ParameterStore configuration.
output "kms_key_id" {
  value = module.env.config_key_id
}

# SQS queue watched for s3 deployment 
# changes when updating the portal ui.
output "s3_change_queue" {
  value = module.env.s3_change_queue
}

