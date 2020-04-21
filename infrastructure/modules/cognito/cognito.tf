variable "signin_url_prefix" {}

variable "metadata_url" {}

variable "sso_redirect_binding_uri" {}

variable "callback_url" {}

variable "logout_url" {}

resource "aws_cognito_user_pool" "pool" {
  name              = "${var.app}-${var.environment}-userpool"
  mfa_configuration = "OFF"
  tags              = var.tags

  admin_create_user_config {
    allow_admin_create_user_only = true
  }
}

resource "aws_cognito_user_pool_domain" "domain" {
  domain       = "${var.signin_url_prefix}-${var.app}-${var.environment}"
  user_pool_id = aws_cognito_user_pool.pool.id
}

output "aws_cognito_sso_url" {
  value = "https://${aws_cognito_user_pool_domain.domain.domain}.auth.${var.region}.amazoncognito.com"
}

output "aws_cognito_return_url" {
  value = "https://${aws_cognito_user_pool_domain.domain.domain}.auth.${var.region}.amazoncognito.com/saml2/idpresponse"
}

output "aws_cognito_audience_restriction" {
  value = "urn:amazon:cognito:sp:${aws_cognito_user_pool.pool.id}"
}


resource "aws_cognito_identity_provider" "provider" {
  user_pool_id  = aws_cognito_user_pool.pool.id
  provider_name = "idp"
  provider_type = "SAML"

  provider_details = {
    MetadataURL           = var.metadata_url
    SSORedirectBindingURI = var.sso_redirect_binding_uri
    IDPSignout            = true
  }

  attribute_mapping = {
    email = "email"
  }
}

resource "aws_cognito_user_pool_client" "client" {
  name                                 = "${var.app}-${var.environment}"
  user_pool_id                         = aws_cognito_user_pool.pool.id
  generate_secret                      = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_scopes                 = ["email", "openid"]
  callback_urls                        = [var.callback_url]
  logout_urls                          = [var.logout_url]
  explicit_auth_flows                  = ["ALLOW_CUSTOM_AUTH", "ALLOW_REFRESH_TOKEN_AUTH", "ALLOW_USER_SRP_AUTH"]
  read_attributes                      = ["address", "birthdate", "email", "email_verified", "family_name", "gender", "given_name", "locale", "middle_name", "name", "nickname", "phone_number", "phone_number_verified", "picture", "preferred_username", "profile", "updated_at", "website", "zoneinfo"]
  supported_identity_providers         = [aws_cognito_identity_provider.provider.provider_name]
  write_attributes                     = ["address", "birthdate", "email", "family_name", "gender", "given_name", "locale", "middle_name", "name", "nickname", "phone_number", "picture", "preferred_username", "profile", "updated_at", "website", "zoneinfo"]
}

output "aws_cognito_signout_url" {
  value = "https://${aws_cognito_user_pool_domain.domain.domain}.auth.us-east-1.amazoncognito.com/logout?client_id=${aws_cognito_user_pool_client.client.id}&logout_uri=${var.logout_url}"
}

output "aws_cognito_signin_url" {
  value = "https://${aws_cognito_user_pool_domain.domain.domain}.auth.us-east-1.amazoncognito.com/oauth2/authorize"
}

output "aws_cognito_token_url" {
  value = "https://${aws_cognito_user_pool_domain.domain.domain}.auth.us-east-1.amazoncognito.com/oauth2/token"
}

output "aws_cognito_client_id" {
  value = aws_cognito_user_pool_client.client.id
}

output "aws_cognito_client_secret" {
  value = aws_cognito_user_pool_client.client.client_secret
}