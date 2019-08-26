# uDeploy #

A simple web based deployment [portal](/docs/PORTAL.md) for AWS resources. 

Deployments can be performed in moments, from anywhere, by any authorized user, at any time.

- Supports: Chrome
- Security: OAuth2

### Features ###

|| Fargate Service | Fargate Task | Lambda Function | S3 Contents ||
|---|---|---|---|---|---|
|User Permissions|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Supports user permissions for editing and deploying environment instances. |
|View Instance Version|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Provides a view into an instance's deployed version details.|
|View Instance Status|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Provides a view into an instances status showing deployments in progress, erroring containers or lambdas, and scaling services and tasks.|
|View Instance Tasks|:white_check_mark:|:white_check_mark:|||Provides a view into what version of a service's tasks are starting or stopping at any given moment.|
|Deploy Version|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Allows drag and drop deployments for all supported AWS resources.|
|Start Instance|:white_check_mark:|:white_check_mark:|:white_check_mark:||Supports quickly starting a stopped service or task.|
|Scale Instance|:white_check_mark:|:white_check_mark:|:white_check_mark:||Supports quickly scaling any service or task.|
|Stop Instance|:white_check_mark:|:white_check_mark:|||Supports quickly  stopping a running service or task.|
|Deployment Notifications|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Sends SNS notifications for AWS resource deployments and other changes in status.|
|Error Notifications|:white_check_mark:|:white_check_mark:|:white_check_mark:||Sends SNS notifications for AWS resource errors. |
|Audit Deployments|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Tracks user deployments via an audit trail.|
|Quick Linking|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Provides direct links to AWS logs, GitHub commits, Jira stories, and many more without browsing through the websites.|
|Deployment Propagation|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Provides option for automatically pushing deployments to QA, UAT, or PROD without user interation keeping environments in sync.|
|Environment Migration|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Supports configuring specific environment variables to be automatically migrated between environments.|
|GitHub Integration|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Provides a quick view into version changes before and after deployments.|


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

* Provide monitoring and deployment portal for AWS resources.
* Provide consitent versioning across different resources.
* Provide high level, easy to understand status notifications for specific environments.
* Provide easily searchable list of project resources, versions, and documentation.
* Provide alternative solution to [Harbor UI](https://github.com/turnerlabs/harbor-ui).
* Provide simple drag-n-drop deployments regardless of the user's AWS experience.
* Provide app and environment level permissions.
* Provide an easy way for integration partners to monitor resources.
* Remove the AWS account permission requirement from deployments and monitoring.

### Non-Goals ###

* Duplicate AWS console functionality.
* Provide continuous integration features.
* Display or modify infrastructure details.

### Tech Stack ###

- Client
    - Bulma (css)
    - Vue.js (javascript)
- Server
    - Echo (go)
- Database
    - MongoDB