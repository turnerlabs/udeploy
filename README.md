# uDeploy #

A simple web based deployment [portal](/docs/PORTAL.md) for AWS resources. 

Authorized users can perform drag-n-drop deployments without understanding the technical aspects of the AWS resource being deployed. Resources like Fargate Tasks, Lambda Functions, and S3 Objects are versioned, deployed, and monitored in the same way reducing complexity.

<details>
  <summary>Reasoning</summary>

### Goals ###

* Expose portal to deploy and monitor AWS resources
* Enable consitent application versioning across multiple AWS resource types
* Secure deployments by application and/or environment
* Allow simple drag-n-drop deployments
* Provide high level environment resource notifications and troubleshooting
* Make projects searchable for quick access to details, versions, and documentation
* Improve resource monitoring for integration partners
* Support authetication with any OAuth2 API
* Enable deployment workflow innovation (avoid third-party timelines)
* Centralize navigation to project resources by linking to scrum boards, config tools, and project documentation
* Assist projects transitioning from [Harbor UI](https://github.com/turnerlabs/harbor-ui)

### Non-Goals ###

* Duplicate AWS console functionality
* Implement continuous integration features
* Display or modify infrastructure

</details>


<details>
  <summary>Features</summary>

|| Fargate Service | Fargate Task | Lambda Function | S3 Contents ||
|---|---|---|---|---|---|
|Authentication|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Supports OAuth2 for authenticating users. |
|Authorization|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Supports user permissions for editing and deploying environment instances. |
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

Supported Browser: Chrome 

</details>

<details>
  <summary>Requirements</summary>

### Instalation ###

|Tool|Version|
|-|-|
|terraform|v0.12|
|docker-compose|v1.24|
|aws cli|v1.16|
|cstore|v3.6|

### Operational ###

|Service|Platform|Purpose|
|-|-|-|
|Route53|AWS|DNS (SSL)|
|ECS Fargate|AWS|Docker Container|
|SQS|AWS|Notifications|
|CloudWatch|AWS|Notifications|
|MongoDB|Atlas|Store|
|OAuth2 Provider|[Azure](docs/OAUTH_AZURE.md)|User Authentication|

</details>

<details>
  <summary>Setup</summary>

1. Create Infrastructure ([Guide](https://github.com/turnerlabs/udeploy-infrastructure))
2. Run portal locally ([Guide](/docs/START.md)) or in Fargate ([Guide](/docs/DEPLOY_PORTAL.md))
  
</details>

<details>
  <summary>Tech Stack</summary>

- Client
    - Bulma (css)
    - Vue.js (javascript)
- Server
    - Echo (go)
- Database
    - MongoDB

</details>


