
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
        "name": "STASH_KEY_PREFIX",
        "value": "${config_subpath}"
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
    "secrets": [],
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