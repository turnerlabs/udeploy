aws_profile = "{{AWS_PROFILE}}"

app         = "udeploy"
environment = "prod"

# zone_name should match the root domain
zone_name = "udeploy.{{ROOT_DOMAIN}}.com"
domain    = "prod.udeploy.{{ROOT_DOMAIN}}.com"

vpc             = "{{VPC}}"
private_subnets = "{{SUBNET_1}},{{SUBNET_2}}"
public_subnets  = "{{SUBNET_3}},{{SUBNET_4}}"
internal        = {{true to use the private_subnets and false to use the public_subnets}}

# Portal configuration acess
saml_role = "{{USER_ROLE}}"
saml_users = [
  "{{USER_EMAIL_1}}",
  "{{USER_EMAIL_2}}",
]

# SSM Parameter Store configuration root path
config_path = "udeploy/infrastructure/portals/prod/.env"

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

linked_account_ids = []
