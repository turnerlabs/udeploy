/*
 * variables.tf
 * Common variables to use in various Terraform files (*.tf)
 */

# The AWS Profile to use
variable "aws_profile" {
}

# The AWS region to use for the dev environment's infrastructure
# Currently, Fargate is only available in `us-east-1`.
variable "region" {
  default = "us-east-1"
}

# Tags for the infrastructure
variable "tags" {
  type = map(string)
}

# The application's name
variable "app" {
}

# The environment that is being built
variable "environment" {
}

# The configuration path in SSM ParameterStore
variable "config_path" {
}

# The port the container will listen on, used for load balancer health check
# Best practice is that this value is higher than 1024 so the container processes
# isn't running at root.
variable "container_port" {
  default = "8080"
}

# The port the load balancer will listen on
variable "lb_port" {
  default = "80"
}

# The load balancer protocol
variable "lb_protocol" {
  default = "HTTP"
}

# Network configuration

# The VPC to use for the Fargate cluster
variable "vpc" {
}

# The private subnets, minimum of 2, that are a part of the VPC(s)
variable "private_subnets" {
}

# The public subnets, minimum of 2, that are a part of the VPC(s)
variable "public_subnets" {
}

variable "zone_name" {
}

variable "record_name" {
}

variable "image" {
  default = "quay.io/turner/udeploy:v0.33.2-rc.18"
}

# Allow other AWS accounts to publish events
# to this account for app status updates
variable "linked_account_ids" {
  type    = list(string)
  default = []
}

variable "ecs_cloudwatch_log_retention_in_days" {
  default = "14"
}

variable "create_lb_http_security_group_rule" {
  default = "1"
}
