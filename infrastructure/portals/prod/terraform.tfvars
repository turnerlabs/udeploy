aws_profile = "{{AWS_PROFILE}}"

app         = "udeploy"
environment = "prod"

# zone_name should match the root domain
zone_name = "udeploy.{{ROOT_DOMAIN}}.com"
domain    = "prod.udeploy.{{ROOT_DOMAIN}}.com"

vpc             = "{{VPC}}"
private_subnets = "{{SUBNET_1}},{{SUBNET_2}}"
public_subnets  = "{{SUBNET_3}},{{SUBNET_4}}"
internal        = {{true (private_subnets) and false (public_subnets)}}

# Portal configuration acess
saml_role = "{{USER_ROLE}}"
saml_users = [
  "{{USER_EMAIL_1}}",
  "{{USER_EMAIL_2}}",
]

# SSM Parameter Store or AWS Secrets Manager configuration root path
config_path = "udeploy/infrastructure/portals/prod/.env"

# Set this value to false to use AWS Secrets Manager
parameter_store = true

tags = {
  application   = "udeploy"
  environment   = "{{ENV}}"
  team          = "{{TEAM}}"
  customer      = "devops"
  contact-email = "{{USER_EMAIL}}"
  product       = "udeploy"
  project       = "udeploy"
  {{OTHER_DESIRED_TAGS}}
}
