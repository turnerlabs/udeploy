# Tags for the infrastructure
variable "tags" {
  type = map(string)
}

# The application's name
variable "app" {
}

# The environment
variable "environment" {
}

# The name of the container to run
variable "container_name" {
  default = "app"
}

# The path to the health check for the load balancer to know if the container(s) are ready
variable "health_check" {
  default = "/ping"
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

# The Hosted Zone for the Route53 records
variable "zone_name" {
}

variable "domain" {
}

# The users (email addresses) from the saml role to give access
# case sensitive
variable "saml_users" {
  type = list(string)
}

variable "saml_role" {
}

variable "config_path" {
}

variable "parameter_store" {
}