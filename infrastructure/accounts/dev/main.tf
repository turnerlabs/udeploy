provider "aws" {
  version = ">= 1.46.0"
  region  = var.region
  profile = var.aws_profile
}

module "dev" {
    source = "../../modules/link"

    region = var.region

    app = var.app
    environment = var.environment

    portal_account_id = var.portal_account_id
}

output "account_id" {
  value = module.dev.account_id
}

output "role_arn" {
  value = module.dev.role_arn
}