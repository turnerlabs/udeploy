
### Run Locally ###

#### 1. Clone ####
```bash
$ git clone git@github.com:turnerlabs/udeploy.git
$ cd udeploy && git checkout v0.28.0-rc
```

### 2. Create Base Infrastucture ####

Replace `{{TOKENS}}` in `./infrastructure/base/terraform.tfvars`.

```bash
$ terraform init -var-file=portals/dev/terraform.tfvars  infrastructure/base 
$ terraform apply -var-file=portals/dev/terraform.tfvars  infrastructure/base
```

#### 3. Configure Portal ####
```bash
$ export ENV=local
$ export URL=http://localhost:8080
```

#### 4. Configure Database ####
Create an empty MongoDB database preferably called `udeploy-dev` on an Atlas M2 (General) cluster or equivalent.

```bash
export DB_URI={{DB_CONNECTION_STRING}}
export DB_NAME={{DB_NAME}}
```

Add an initial admin user to the `users` collection. Additional users can be added through the portal.

```
db.users.insert({"admin":true,"email":"User.Email@domain.com","apps":{}})
```

IMPORTANT: The email address is case sensitive.

#### 5. Configure Event Notifications ####
```bash
$ export SQS_CHANGE_QUEUE=udeploy-local-notification-queue.fifo
$ export SQS_ALARM_QUEUE=udeploy-local-alarm-queue
$ export SQS_S3_QUEUE=udeploy-local-s3-queue
$ export SNS_ALARM_TOPIC_ARN=arn:aws:sns:us-east-1:{{ACCOUNT_ID}}:{{SNS_ALARM_TOPIC_NAME}}
```

#### 6. Configure User Authentication ####
The portal uses oauth2 for authenticating users before loading database authorization details. To configure an (OIDC) OpenID Provider, set the following environment variables and ensure the portal is configured with the desired provider.

This example uses the Azure provider. [Register](OAUTH_AZURE.md) the portal with Azure to generate the tokens below.

```bash
$ export OAUTH_CLIENT_ID={{CLIENT_ID}} 
$ export OAUTH_CLIENT_SECRET={{CLIENT_SECRET}}
$ export OAUTH_AUTH_URL=https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/v2.0/authorize)
$ export OAUTH_TOKEN_URL=https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/v2.0/token)
$ export OAUTH_REDIRECT_URL=http://localhost:8080/oauth2/response)
$ export OAUTH_SIGN_OUT_URL=https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/logout?client_id={{CLIENT_ID}})
$ export OAUTH_SESSION_SIGN=F1Li3rvehgcrF8DMHJ7OyxO4w9Y3D
```

The `OAUTH_SESSION_SIGN` can be any secure string.

#### 7. Configure Console ####
Browser quick link to the AWS console.
```bash
$ export CONSOLE_LINK=https://aws.amazon.com/
```

#### 8. Configure In-Memory Cache ####
Starts caching applications when the container starts to improve perfomance.
```bash
$ export PRE_CACHE=true
```

#### 9. Run #### 

```bash
$ go run . # Go 1.13.x or higher required
```

PORTAL: http://localhost:8080