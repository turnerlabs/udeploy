app = "udeploy"

aws_profile = "{{AWS_PROFILE}}"

# Access to ECR
saml_role = "{{USER_ROLE}}" //the role that will be used to access the tf remote state s3 bucket and the ecr images
saml_user = "{{USER_EMAIL}}"

# Base domain for all environments
domain = "udeploy.{{ROOT_DOMAIN}}.com"

# Uncomment to point the A record to the prod instance.
# alias_zone_id = "{{ALIAS_ZONE_ID}}"
# alias_name = "{{ALIAS_NAME}}"

tags = {
  application      = "udeploy"
  environment      = "{{ENV}}"
  team             = "{{TEAM}}"
  customer         = "devops"
  contact-email    = "{{USER_EMAIL}}" // user or team email
  product          = "udeploy"
  project          = "udeploy"
}