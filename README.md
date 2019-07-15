# uDeploy #

A simple web based deployment [portal](/docs/PORTAL.md) for AWS resources. 

Deployments can be performed in moments, from anywhere, by any authorized user, at any time.

### Features ###

Supports: Chrome, Security: OAuth2

|| Fargate Service | Fargate Task | Lambda Function | S3 Contents |
|---|---|---|---|---|
|User Permissions|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|
|View Instance Version|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|
|View Instance Status|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|
|View Instance Tasks|:white_check_mark:|:white_check_mark:|||
|Deploy Version|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|
|Start Instance|:white_check_mark:|:white_check_mark:|:white_check_mark:||
|Scale Instance|:white_check_mark:|:white_check_mark:|:white_check_mark:||
|Stop Instance|:white_check_mark:|:white_check_mark:|||
|Deployment Notifications|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|
|Error Notifications|:white_check_mark:|:white_check_mark:|:white_check_mark:||
|Audit Deployments|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|
|Quick Linking|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|
|AWS Log Links|:white_check_mark:|:white_check_mark:|:white_check_mark:||
|Deployment Propagation|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|
|Environment Migration|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|
|GitHub Integration|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|


### Get Started ###

[Run](/docs/START.md) container from scratch.

### Hosting Solution ###

uDeploy currently services resources within a single AWS account.

1. Create [Infrastructure](https://github.com/turnercode/ams-udeploy-infrastructure)
    - DNS + SSL (Route53)
    - Queues (SQS)
    - Alerts (CloudWatch)
    - MongoDB Database (Atlas)
    - OAuth2 (Azure)

2. Deploy [Service](/docs/START.md)
    - Application Configuration
    - Fargate Service (AWS)

### Goals ###

* Provide a deployment portal for simple AWS service, task, and lambda deployments.
* Provide notifications for instance/environment statuses.

### Non-Goals ###

* Duplicate the AWS console.
* Display or modify infrastructure details.

### Tech Stack ###

- Client
    - Bulma (css)
    - Vue.js (javascript)
- Server
    - Echo (go)
- Database
    - MongoDB

