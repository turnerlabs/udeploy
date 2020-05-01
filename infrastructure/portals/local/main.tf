terraform {
  required_version = ">= 0.12.0"

  backend "s3" {
    region  = "us-east-1"
    profile = "{{PROFILE}}"
    bucket  = "tf-state-udeploy"
    key     = "local.terraform.tfstate"
  }
}

# The AWS Profile to use
variable "aws_profile" {
}

provider "aws" {
  version = ">= 1.46.0"
  region  = var.region
  profile = var.aws_profile
}

# output

# Command to set the AWS_PROFILE
output "aws_profile" {
  value = var.aws_profile
}

data "aws_caller_identity" "current" {
}

