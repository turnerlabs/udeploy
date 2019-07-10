
### Get Started ###

#### Clone ####
```bash
$ git clone git@github.com:turnerlabs/udeploy.git && cd udeploy
```

#### Configure Database ####
```bash
$ export DB_URI=mongodb://localhost:27017 
$ export DB_NAME=udeploy
```

In a `users` database collection, add an entry for the initial admin user to get started. The rest of the users can be added through the portal by the admin user. 

The email address is case sensitive.

Example:
```json
{
    "admin": true,
    "email": "user.email@somewhere.com",
    "apps": {
        "app-name": {
            "claims": ["edit"]
        }
    }
}
```

If an application is not listed under the `apps` key in the user policy, the user will not have access to view the application. 

If an application is listed under the `apps` key in the user policy, the user will be able to view all the applications instances, but not be able to `scale`, `deploy`, or `edit` the instances without the specific claims.

Only non-admin users need `edit` claims for making modifications, but any users, including admins, who intend to perform deployments need `scale` and `deploy` claims specified.

#### Configure Event Notifications ####
```bash
$ export ENV=local
$ export SQS_CHANGE_QUEUE=udeploy-local-notification-queue.fifo
$ export SQS_ALARM_QUEUE=udeploy-local-alarm-queue
$ export SQS_S3_QUEUE=udeploy-local-s3-queue
$ export SNS_ALARM_TOPIC_ARN=arn:aws:sns:us-east-1:{{ACCOUNT_ID}}:{{SNS_ALARM_TOPIC_NAME}
$ export URL=http://localhost:8080
```

#### Configure User Authentication ####
uDeploy uses oauth2 for authenticating users before loading authorization details from the database. To configure an (OIDC) OpenID Provider, set the following environment variables and ensure uDeploy is configured with the provider.

```bash
$ export OAUTH_CLIENT_ID=XXXXXXX 
$ export OAUTH_CLIENT_SECRET=XXXXXXX
$ export OAUTH_AUTH_URL=XXXXXXX (i.e. https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/v2.0/authorize)
$ export OAUTH_TOKEN_URL=XXXXXXX (i.e. https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/v2.0/token)
$ export OAUTH_REDIRECT_URL=XXXXXXX (i.e. http://localhost:8080/oauth2/response)
$ export OAUTH_SIGN_OUT_URL=XXXXXXX (i.e. https://login.microsoftonline.com/{{TENANT_ID}}/oauth2/logout?client_id={{CLIENT_ID}})
$ export OAUTH_SESS_SIGN=XXXXXXX (i.e. F1Li3rvehgcrF8DMHJ7OyxO4w9Y3D)
```

#### Configure Console ####
Browser quick link to the AWS console.
```bash
$ export CONSOLE_LINK=https://aws.amazon.com/
```


#### Start #### 

```bash
$ go run . # Go 1.11.x or higher required
```

PORTAL: http://localhost:8080