// only lower-case letters, numbers, and hyphens
signin_url_prefix = "{{AWS_COGNITO_SIGNIN_URL_PREFIX}}"

metadata_url              = "{{OKTA_METADATA_URL}}"
sso_redirect_binding_uri  = "{{OKTA_SSO_REDIRECT_BINDING_URI}}"

callback_url         = "https://{{DOMAIN}}/oauth2/response"
logout_url           = "https://{{OKTA_DOMAIN}}/logout.aspx?AppID={{APP_ID}}"