aws_profile = "{{AWS_PROFILE}}"

app = "udeploy"
environment = "{{ENV}}"

# zone_name should match the root domain
zone_name = "udeploy.{{ROOT_DOMAIN}}.com"
domain = "{{ENV}}.udeploy.{{ROOT_DOMAIN}}.com"

vpc = "{{VPC}}"
private_subnets = "{{SUBNET_1}},{{SUBNET_2}}"
public_subnets = "{{SUBNET_3}},{{SUBNET_4}}"

# Portal configuration acess
saml_role = "{{USER_ROLE}}"
saml_users = [
  "{{USER_EMAIL_1}}",
  "{{USER_EMAIL_2}}",
]

# ECR Image
image = "{{IMAGE}}"

# SSM Parameter Store configuration root path
config_path = "{{CONFIG_PATH}}"

tags = {
  application      = "udeploy"
  environment      = "{{ENV}}"
  team             = "{{TEAM}}"
  customer         = "devops"
  contact-email    = "{{USER_EMAIL}}"
  product          = "udeploy"
  project          = "udeploy"
}