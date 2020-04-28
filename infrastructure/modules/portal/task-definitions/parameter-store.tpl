[
    {
      "name": "${container_name}",
      "image": "${image}",
      "essential": true,
      "portMappings": [
        {
          "protocol": "tcp",
          "containerPort": ${container_port},
          "hostPort": ${container_port}
        }
      ],
      "environment": [
        {
          "name": "PORT",
          "value": "${container_port}"
        },
        {
          "name": "HEALTHCHECK",
          "value": "${health_check}"
        },
        {
          "name": "ENV",
          "value": "${environment}"
        },
        {
          "name": "PRE_CACHE",
          "value": "true"
        },
        {
          "name": "URL",
          "value": "https://${record_name}"
        },
        {
          "name": "SQS_CHANGE_QUEUE",
          "value": "${sqs_change_queue}"
        },
        {
          "name": "SQS_ALARM_QUEUE",
          "value": "${sqs_alarm_queue}"
        },
        {
          "name": "SQS_S3_QUEUE",
          "value": "${sqs_s3_queue}"
        },
        {
          "name": "SNS_ALARM_TOPIC_ARN",
          "value": "${sns_alarm_topic_arn}"
        }
      ],
      "secrets": [
        {
          "valueFrom": "/${config_path}/DB_URI",
          "name": "DB_URI"
        },
        {
          "valueFrom": "/${config_path}/OAUTH_SIGN_OUT_URL",
          "name": "OAUTH_SIGN_OUT_URL"
        },
        {
          "valueFrom": "/${config_path}/CONSOLE_LINK",
          "name": "CONSOLE_LINK"
        },
        {
          "valueFrom": "/${config_path}/DB_NAME",
          "name": "DB_NAME"
        },
        {
          "valueFrom": "/${config_path}/OAUTH_CLIENT_ID",
          "name": "OAUTH_CLIENT_ID"
        },
        {
          "valueFrom": "/${config_path}/OAUTH_CLIENT_SECRET",
          "name": "OAUTH_CLIENT_SECRET"
        },
        {
          "valueFrom": "/${config_path}/OAUTH_SESSION_SIGN",
          "name": "OAUTH_SESSION_SIGN"
        },
        {
          "valueFrom": "/${config_path}/OAUTH_TOKEN_URL",
          "name": "OAUTH_TOKEN_URL"
        },
        {
          "valueFrom": "/${config_path}/OAUTH_AUTH_URL",
          "name": "OAUTH_AUTH_URL"
        },
        {
          "valueFrom": "/${config_path}/OAUTH_REDIRECT_URL",
          "name": "OAUTH_REDIRECT_URL"
        },
        {
          "valueFrom": "/${config_path}/OAUTH_SCOPES",
          "name": "OAUTH_SCOPES"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/fargate/service/${app}-${environment}",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]