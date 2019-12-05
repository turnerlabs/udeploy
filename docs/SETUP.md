
## Deploy Portal to Fargate ##


### 1. Clone Infrastructure ###
```bash
$ mkdir udeploy && cd udeploy
$ git clone --branch v0.29.2-rc git@github.com:turnerlabs/udeploy.git infrastructure
```

### 2. Create Base Infrastucture ###

Optionally, customize [infrastructure](BASE.md).

Replace `{{TOKENS}}` in [infrastructure/base/terraform.tfvars](/infrastructure/base/terraform.tfvars).
```bash
$ terraform init -var-file=infrastructure/base/terraform.tfvars  infrastructure/base 
$ terraform apply -var-file=infrastructure/base/terraform.tfvars  infrastructure/base
```

### 3. Initilize Configuration ###

```bash
$ cp infrastructure/env.template infrastructure/portals/prod/.env
```

In `infrastructure/portals/prod/.env` Replace `{{DOMAIN}}` with the portal domain.

```
URL=https://{{DOMAIN}}
```

### 4. Create Database ###

Create an empty MongoDB database preferably called `udeploy-prod` on an Atlas M2 (General) cluster or equivalent. [Want to Terraform Atlas MongoDB?](ATLAS.md) If not, delete [atlas.tf](/infrastructure/portals/prod/atlas.tf).

Replace `{{TOKENS}}` in `./infrastructure/portals/prod/.env` file.
```
DB_URI={{DB_CONNECTION_STRING}}
DB_NAME={{DB_NAME}}
```

Add an initial admin user to the `users` collection. Additional users can be added through the portal.

```
use {{DB_NAME}}
db.users.insert({"admin":true,"email":"User.Email@domain.com","apps":{}})
```

IMPORTANT: The email address is case sensitive.

### 5. Configure User Authentication ###

The portal uses oauth2 for authenticating users before loading user authorization.

This example uses the Azure OpenID provider. 

After [registering](OAUTH_AZURE.md) the portal with Azure, replace `{{TOKENS}}` in `./infrastructure/portals/prod/.env` file.

```bash
OAUTH_CLIENT_ID={{CLIENT_ID}} 
OAUTH_CLIENT_SECRET={{CLIENT_SECRET}}
OAUTH_AUTH_URL=https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/v2.0/authorize
OAUTH_TOKEN_URL=https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/v2.0/token
OAUTH_REDIRECT_URL=https://{{DOMAIN}}/oauth2/response
OAUTH_SIGN_OUT_URL=https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/logout?client_id={{CLIENT_ID}}
OAUTH_SESSION_SIGN=F1Li3rvehgcrF8DMHJ7OyxO4w9Y3D
```

The `OAUTH_SESSION_SIGN` should be updated to a any secure string.

### 6. Configure Event Notifications ###

Replace `{{TOKENS}}` in `./infrastructure/portals/prod/.env` file.

```bash
SQS_CHANGE_QUEUE=udeploy-prod-notification-queue.fifo
SQS_ALARM_QUEUE=udeploy-prod-alarm-queue
SQS_S3_QUEUE=udeploy-prod-s3-queue
SNS_ALARM_TOPIC_ARN=arn:aws:sns:us-east-1:{{ACCOUNT_ID}}:{{APP}}-prod-alarms
```

### 7. Create Portal Infrastucture ####

Replace `{{TOKENS}}` in [infrastructure/portals/prod/terraform.tfvars](/infrastructure/portals/prod/terraform.tfvars).

```bash
$ terraform init -var-file=infrastructure/portals/prod/terraform.tfvars infrastructure/portals/prod
$ terraform apply -var-file=infrastructure/portals/prod/terraform.tfvars infrastructure/portals/prod
```

### 8. Push Configuration to Parameter Store ### 

Install [cstore](https://github.com/turnerlabs/cstore) and run the following commands from the repository root to store configuration in SSM Parameter Store.

```bash
$ export AWS_REGION=us-east-1
$ export AWS_PROFILE=aws-account-profile
```

```bash
$ cstore push infrastructure/portals/prod/.env -s aws-parameter -t prod
```

When prompted, set context to `udeploy` and the KMS Key ID to the `kms_key_id` from the Terraform output.

 ### 9. Link Other AWS Accounts (optional) ### 

 To deploy resources accross multiple AWS accounts, provide permissions to each additional AWS account the portal should control. 

 Duplicate the folder `infrastructure/accounts/dev` for each account `infrastructure/accounts/{{ACCOUNT_IDENTIFIER}}` and following the intructions.

 Replace `{{TOKENS}}` in `infrastructure/accounts/{{ACCOUNT_IDENTIFIER}}/terraform.tfvars`.
```bash
$ terraform init -var-file=infrastructure/accounts/{{ACCOUNT_IDENTIFIER}}/terraform.tfvars  infrastructure/accounts/{{ACCOUNT_IDENTIFIER}} 
$ terraform apply -var-file=infrastructure/accounts/{{ACCOUNT_IDENTIFIER}}/terraform.tfvars  infrastructure/accounts/{{ACCOUNT_IDENTIFIER}}
```

Update `linked_account_ids` in [infrastructure/base/terraform.tfvars](/infrastructure/base/terraform.tfvars) with account ids of all linked accounts.

```
$ terraform apply -var-file=infrastructure/base/terraform.tfvars  infrastructure/base
```



