
## Deploy Portal to Fargate ##


### 1. Clone Repository ###
```bash
$ git clone git@github.com:turnerlabs/udeploy.git
$ cd udeploy && git checkout v0.28.0-rc
```

The following commands should be executed from the repository root unless otherwise specfied.

### 2. Create Base Infrastucture ###

These are the files included and what they accomplish:

- [./infrastructure/base/ecr.tf](infrastructure/base/ecr.tf) (optional)
    ```
    - Creates an AWS ECR and the permissions necessary. 
    - Keep this file if this is the location in which you intend to store the docker images.
    ```
- [./infrastructure/base/main.tf](infrastructure/base/main.tf) (required)
    `- Main configuration`
- [./infrastructure/base/route53.tf](infrastructure/base/route53.tf) (optional)
    ```
    - Creates an AWS Route53 Hosted Zone, a Certificate Manager, and the DNS Records necessary to validate the Certificate. 
    - It will also validate the certificate and wait for the full validation to be completed. 
    - For the validation to fully happen, the Hosted Zone has to propagate and be available in the internet.
    - Form more information check this link: https://docs.aws.amazon.com/acm/latest/userguide/gs-acm-validate-dns.html
    ```
- [./infrastructure/base/state.tf](infrastructure/base/state.tf) (optional)
    ```
    - Creates an s3 bucket to store the tfstate files that will be generated when the terraform apply command runs.
    - Although not required, these approach is highly recommended to avoid having your state files being public in your repo.
    - An alternative is to use another s3 bucket where state is already kept for your account/projects
    ```
- [./infrastructure/base/terraform.tfvars](infrastructure/base/terraform.tfvars) (required)
    `File containing the main configuration for the group of files in the base folder`
- [./infrastructure/base/variables.tf](infrastructure/base/variables.tf) (required)
    `terraform declaration of variables`
- [./infrastructure/base/versions.tf](infrastructure/base/versions.tf) (required)
    `terraform version requirements`

1. Remove any files that are not necessary for your use case. Then run the following commands
1. The following steps are only necessary if any of the optional files above is still left. Otherwise the base folder will not create any resources.
1. Replace `{{TOKENS}}` in [./infrastructure/base/ecr.tf](infrastructure/terraform.tfvars).
1. 
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

Create an empty MongoDB database preferably called `udeploy-prod` on an Atlas M2 (General) cluster or equivalent.

Replace `{{TOKENS}}` in `./infrastructure/portals/prod/.env` file.
```
DB_URI={{DB_CONNECTION_STRING}}
DB_NAME={{DB_NAME}}
```

Add an initial admin user to the `users` collection. Additional users can be added through the portal.

```
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

### 6. Configure Events ###

Replace `{{TOKENS}}` in `./infrastructure/portals/prod/.env` file.

```bash
SQS_CHANGE_QUEUE=udeploy-prod-notification-queue.fifo
SQS_ALARM_QUEUE=udeploy-prod-alarm-queue
SQS_S3_QUEUE=udeploy-prod-s3-queue
SNS_ALARM_TOPIC_ARN=arn:aws:sns:us-east-1:{{ACCOUNT_ID}}:{{APP}}-prod-alarms
```

### 7. Specify/Build Docker Image ####

Update `./infrastructure/portals/prod/terraform.tfvars` with the image.

```
image = "{{IMAGE}}"
```

If a image does not exist, create one before updating the `terraform.tfvars` file using the following directions otherwise skip to step `#8`.

```bash
$ export AWS_DEFAULT_REGION=us-east-1
$ export AWS_PROFILE=aws-account-profile
```

```bash
$ login=$(aws ecr get-login --no-include-email) && eval "$login"
```

```bash
$ docker build -t {{ACCOUNT_ID}}.dkr.ecr.us-east-1.amazonaws.com/{{APP}}:v0.28.0-rc --build-arg version=v0.28.0-rc.1 .
$ docker push {{ACCOUNT_ID}}.dkr.ecr.us-east-1.amazonaws.com/{{APP}}:v0.28.0-rc
```

### 8. Create Portal Infrastucture ####

Replace `{{TOKENS}}` in `./infrastructure/portals/prodterraform.tfvars`.

```bash
$ terraform init -var-file=portals/prod/terraform.tfvars infrastructure/portals/prod
$ terraform apply -var-file=portals/prod/terraform.tfvars infrastructure/portals/prod
```

To point the A record to the prod load balancer, copy the previous command output values `alias_zone_id` and `alias_name` to `./infrastructure/base/terraform.tfvars` and uncomment the lines related to the alias in these two files.

* `./infrastructure/base/route53.tf`
* `./infrastructure/base/terraform.tfvars`

```bash
$ terraform apply -var-file=portals/prod/terraform.tfvars infrastructure/portals/prod
```


### 9. Push Configuration to Parameter Store ### 

Install [cstore](https://github.com/turnerlabs/cstore) and run the following commands from the repository root to store configuration in SSM Parameter Store.

```bash
$ export AWS_REGION=us-east-1
$ export AWS_PROFILE=aws-account-profile
```

```bash
$ cstore push infrastructure/portals/prod/.env -s aws-parameter -t prod
```

When prompted, set context to `udeploy` and the KMS Key ID to the `kms_key_id` from the Terraform output.

 