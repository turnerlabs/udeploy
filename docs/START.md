
### Run Locally ###

#### Clone ####
```bash
$ git clone git@github.com:turnerlabs/udeploy.git && cd udeploy
```

#### Configure Portal ####
```bash
$ export ENV=local
$ export URL=http://localhost:8080
```

#### Configure Database ####
```bash
$ export DB_URI=mongodb://localhost:27017 
$ export DB_NAME=udeploy
```

In a `users` database collection, add an entry for the initial admin user to get started. The rest of the users can be added through the portal by the admin user. 

The email address is case sensitive.

```
db.users.insert({"admin":true,"email":"user.email@somewhere.com","apps":{}})
```

#### Configure Event Notifications ####
```bash
$ export SQS_CHANGE_QUEUE=udeploy-local-notification-queue.fifo
$ export SQS_ALARM_QUEUE=udeploy-local-alarm-queue
$ export SQS_S3_QUEUE=udeploy-local-s3-queue
$ export SNS_ALARM_TOPIC_ARN=arn:aws:sns:us-east-1:{{ACCOUNT_ID}}:{{SNS_ALARM_TOPIC_NAME}}
```

#### Configure User Authentication ####
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

#### Configure Console ####
Browser quick link to the AWS console.
```bash
$ export CONSOLE_LINK=https://aws.amazon.com/
```

#### Configure In-Memory Cache ####
Starts caching applications when the container starts to improve perfomance.
```bash
$ export PRE_CACHE=true
```


#### Run #### 

```bash
$ go run . # Go 1.13.x or higher required
```

PORTAL: http://localhost:8080