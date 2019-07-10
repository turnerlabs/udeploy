
### Configure Applications and Instances ###

|Field|Required|Describes|
|-----|-----|-----|
|name|yes| What the application will be called in the ui. |
|type|yes| How the application should be treated. As a `service`, `scheduled-task`, or `lambda` in AWS. |
|instance.{NAME}|yes| What the instance/environment will be called in the ui. |
|instance.{NAME}.order|yes| How to order the instances in the ui. (i.e. dev, qa, uat, prod). |
|instance.{NAME}.cluster|`scheduled-task`,`service`| What AWS cluster the instance is categorized under. |
|instance.{NAME}.service|`scheduled-task`,`service`| What AWS service the instance is categorized under. |
|instance.{NAME}.eventRule|`scheduled-task`| Which event rule the instance should use. (only applies to `scheduled-tasks`)|
|instance.{NAME}.functionName|`lambda`| Which lambda function should displayed. (only applies to `lambda`)|
|instance.{NAME}.functionAlias|`lambda`| Which alias to update with the new version during deployments. (only applies to `lambda`)|
|instance.{NAME}.repository|`scheduled-task`,`service`| Which ECR repository the images should be pulled from when populating the available versions list. |
|instance.{NAME}.deployCode|| What code a user is required to enter when when deployments and other actions are requested. |
|instance.{NAME}.taskDefinition.family|`scheduled-task`,`service`| Which task definition family to use when listing versions all the user can deploy. |
|instance.{NAME}.taskDefinition.imageTagEx|yes| How to parse the image's tag or the lambda function description to get a version. The first regex capturing group represents the version and the second, when present, represents the build number. |
|instance.{NAME}.taskDefinition.cloneEnvVars|| Which task definition environment variables to copy from the source instance to the target instance when deploying. If empty, none of the source's environment variables will be copied. |
|instance.{NAME}.taskDefinition.registry|| Which instance should be the default source for populating the available versions list. |
|instance.{NAME}.taskDefinition.revisions|`scheduled-task`,`service`| How many task definition revisions starting from the most recent should be considered when building a list of application versions. |