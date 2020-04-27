# uDeploy #

A simple web based deployment [portal](/docs/PORTAL.md) for AWS resources. 

Authorized users can perform drag-n-drop deployments without understanding the technical aspects of the AWS resource being deployed. Resources like Fargate Tasks, Lambda Functions, and S3 Objects are versioned, deployed, and monitored in the same way reducing complexity.

This portal compliments CI/CD processes like GitHub Actions, Jenkins, TeamCity, AWS CodeBuild, CircleCI, and the like with a user friendly GUI. It is not intended to be a replacement for CI/CD pipelines.

<details>
  <summary>Why?</summary>

### Goals ###

* Create portal to deploy and monitor AWS resources
* Make projects searchable for quick access to details, versions, and documentation
* Enable consitent application versioning across multiple AWS resource types
* Secure deployments by application and/or environment
* Allow simple drag-n-drop deployments
* Provide high level environment resource notifications and troubleshooting
* Improve resource monitoring and interaction for integration partners
* Support authentication with any OAuth2 API
* Enable deployment workflow innovation (avoid third-party timelines and costs)
* Centralize navigation to project resources by linking to scrum boards, config tools, and project documentation
* Assist projects transitioning from [Harbor UI](https://github.com/turnerlabs/harbor-ui)

### Non-Goals ###

* Duplicate AWS console functionality
* Implement continuous integration features
* Display or modify infrastructure
* Replace existing CI/CD piplines

</details>


<details>
  <summary>Features</summary>

|| Fargate Service | Fargate Task | Lambda Function | S3 Contents ||
|---|---|---|---|---|---|
|Multiple Accounts|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Manages AWS resources accross multiple accounts with a single portal. |
|Authentication|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Supports OAuth2 for authenticating users. |
|Authorization|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Supports role based user permissions for editing and deploying environment instances. |
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
|Deployment Propagation|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Provides option for automatically pushing deployments to QA, UAT, or PROD without user interaction keeping environments in sync.|
|Environment Migration|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Supports configuring specific environment variables to be automatically migrated between environments.|
|GitHub Integration|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Provides a quick view into version changes before and after deployments.|
|Projects|:white_check_mark:|:white_check_mark:|:white_check_mark:|:white_check_mark:|Allows applications to be grouped by projects for organization and function.|

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
  <summary>Infrastructure</summary>

* [Local](/docs/START.md)
* [Fargate](/docs/SETUP.md)
  
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

<details>
  <summary>Architecture</summary>

The portal functions as a passive monitoring and active deployment platform. The architecture diagram is divided into two sections, monitoring and user actions.

![uDeploy Architecture](/image/architecture.png "uDeploy Architecture")

</details>

<details>
  <summary>Publish Build</summary>

NOTE: semver tag format required `v1.0.0-rc`
```bash
$ git tag {{TAG}}
$ git push origin {{TAG}}
```
Once the GitHub Actions build is complete the Docker image tagged {{TAG}} release will be published to Quay.io. 
</details>

<details>
  <summary>Configure Apps</summary>

Step by step instructions for configuring applications in a running instance of the portal.

* [Lambda Functions](/docs/guides/lambda/README.md)
* [Fargate Tasks](/docs/guides/fargate/README.md)
* [S3 Buckets](/docs/guides/s3/README.md)

</details>