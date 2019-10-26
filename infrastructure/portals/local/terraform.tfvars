aws_profile = "{{AWS_PROFILE}}"
region = "us-east-1"

app = "udeploy"
environment = "{{ENV}}"

vpc = "{{VPC}}"
private_subnets = "{{SUBNET_1}},{{SUBNET_2}}"
public_subnets = "{{SUBNET_3}},{{SUBNET_4}}"

tags = {
  application      = "udeploy"
  environment      = "{{ENV}}"
  team             = "{{TEAM}}"
  customer         = "devops"
  contact-email    = "{{USER_EMAIL}}"
  product          = "udeploy"
  project          = "udeploy"
}