aws_profile = "{{AWS_PROFILE}}"
region = "us-east-1"

app = "udeploy"
environment = "{{ENV}}"

tags = {
  application      = "udeploy"
  environment      = "{{ENV}}"
  team             = "{{TEAM}}"
  customer         = "devops"
  contact-email    = "{{USER_EMAIL}}"
  product          = "udeploy"
  project          = "udeploy"
}