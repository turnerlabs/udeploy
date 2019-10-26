
### Deploy Portal to Fargate ###

#### Clone Repository ####
```bash
$ git clone git@github.com:turnerlabs/udeploy.git
$ git checkout v0.28.0-rc
```

#### Create Configuration File ####

Create the following `.env` file at the repository root.

```
ENV={{ENV}}
URL=https://{{DOMAIN}}
CONSOLE_LINK=https://aws.amazon.com/
PRE_CACHE=true
```

#### Create Database ####

Create a MongoDB database.

Add the database configuration to the `.env` file.
```
DB_URI={{DB_CONNECTION_STRING}}
DB_NAME={{APP}}
```

In a `users` collection, add an entry for the initial admin user. The rest of the users can be added through the portal by the admin user. 

The email address is case sensitive.

```
db.users.insert({"admin":true,"email":"user.email@somewhere.com","apps":{}})
```

#### Configure Event Notifications ####

Add the event configuration to the `.env` file.

```bash
SQS_CHANGE_QUEUE=udeploy-{{ENV}}-notification-queue.fifo
SQS_ALARM_QUEUE=udeploy-{{ENV}}-alarm-queue
SQS_S3_QUEUE=udeploy-{{ENV}}-s3-queue
SNS_ALARM_TOPIC_ARN=arn:aws:sns:us-east-1:{{ACCOUNT_ID}}:{{APP}}-{{ENV}}-alarms
```

#### Configure User Authentication ####

The portal uses oauth2 for authenticating users before loading database authorization details. To configure an (OIDC) OpenID Provider, set the following environment variables and ensure the portal is configured with the desired provider.

This example uses the Azure provider. 

After [registering](OAUTH_AZURE.md) the portal with Azure add the auth configuration to the `.env` file.

```bash
OAUTH_CLIENT_ID={{CLIENT_ID}} 
OAUTH_CLIENT_SECRET={{CLIENT_SECRET}}
OAUTH_AUTH_URL=https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/v2.0/authorize
OAUTH_TOKEN_URL=https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/v2.0/token
OAUTH_REDIRECT_URL=https://{{DOMAIN}}/oauth2/response
OAUTH_SIGN_OUT_URL=https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/logout?client_id={{CLIENT_ID}}
OAUTH_SESSION_SIGN=F1Li3rvehgcrF8DMHJ7OyxO4w9Y3D
```

The `OAUTH_SESSION_SIGN` can be any secure string.

#### Push Configuration to Parameter Store #### 

Install [cstore](https://github.com/turnerlabs/cstore) and run the following commands from the repository root to store configuration in SSM Parameter Store.

```bash
$ export AWS_REGION=us-east-1
$ export AWS_PROFILE={{PROFILE}} 
```

```bash
$ cstore push .env -s aws-parameter -t {{ENV}
```

#### Push Image to ECR #### 

Run commands from repository root.

```bash
$ export AWS_DEFAULT_REGION=us-east-1
$ export AWS_PROFILE={{PROFILE}} 
```

```bash
$ login=$(aws ecr get-login --no-include-email) && eval "$login"
```

```bash
$ docker build -t {{ACCOUNT_ID}}.dkr.ecr.us-east-1.amazonaws.com/{{APP}}:v0.28.0-rc --build-arg version=v0.28.0-rc.1 .
$ docker push {{ACCOUNT_ID}}.dkr.ecr.us-east-1.amazonaws.com/{{APP}}:v0.28.0-rc
```

#### Deploy Portal #### 

Add `{{ACCOUNT_ID}}.dkr.ecr.us-east-1.amazonaws.com/{{APP}}:v0.28.0-rc` image to [Task Definition](https://github.com/turnerlabs/udeploy-infrastructure) and apply Terraform changes.
 