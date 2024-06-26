/**
 * Elastic Container Service (ecs)
 * This component is required to create the Fargate ECS service. It will create a Fargate cluster
 * based on the application name and enironment. It will create a "Task Definition", which is required
 * to run a Docker container, https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_definitions.html.
 * Next it creates a ECS Service, https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs_services.html
 * It attaches the Load Balancer created in `lb.tf` to the service, and sets up the networking required.
 * It also creates a role with the correct permissions. And lastly, ensures that logs are captured in CloudWatch.
 *
 * When building for the first time, it will install a "default backend", which is a simple web service that just
 * responds with a HTTP 200 OK. It's important to uncomment the lines noted below after you have successfully
 * migrated the real application containers to the task definition.
 */

# How many containers to run
variable "replicas" {
  default = "1"
}

# The name of the container to run
variable "container_name" {
  default = "app"
}

# The minimum number of containers that should be running.
# Must be at least 1.
# used by both autoscale-perf.tf and autoscale.time.tf
# For production, consider using at least "2".
variable "ecs_autoscale_min_instances" {
  default = "1"
}

# The maximum number of containers that should be running.
# used by both autoscale-perf.tf and autoscale.time.tf
variable "ecs_autoscale_max_instances" {
  default = "8"
}

resource "aws_ecs_cluster" "app" {
  name = "${var.app}-${var.environment}"
  tags = var.tags

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}

resource "aws_ecs_task_definition" "app" {
  family                   = "${var.app}-${var.environment}"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecsTaskExecutionRole.arn

  # defined in role.tf
  task_role_arn = aws_iam_role.app_role.arn

  container_definitions = <<DEFINITION
[
  {
    "name": "${var.container_name}",
    "image": "${var.image}",
    "essential": true,
    "portMappings": [
      {
        "protocol": "tcp",
        "containerPort": ${var.container_port},
        "hostPort": ${var.container_port}
      }
    ],
    "environment": [
      {
        "name": "PORT",
        "value": "${var.container_port}"
      },
      {
        "name": "HEALTHCHECK",
        "value": "${var.health_check}"
      },
      {
        "name": "APP",
        "value": "${var.app}"
      },
      {
        "name": "ENVIRONMENT",
        "value": "${var.environment}"
      },
      {
        "name": "PRE_CACHE",
        "value": "true"
      },
      {
        "name": "URL",
        "value": "https://${var.record_name}"
      },
      {
        "name": "KMS_KEY_ID",
        "value": "${aws_kms_key.config.id}"
      }
    ],
    "secrets": [
      {
        "valueFrom": "/${var.config_path}/DB_URI",
        "name": "DB_URI"
      },
      {
        "valueFrom": "/${var.config_path}/OAUTH_SIGN_OUT_URL",
        "name": "OAUTH_SIGN_OUT_URL"
      },
      {
        "valueFrom": "/${var.config_path}/CONSOLE_LINK",
        "name": "CONSOLE_LINK"
      },
      {
        "valueFrom": "/${var.config_path}/ENV",
        "name": "ENV"
      },
      {
        "valueFrom": "/${var.config_path}/SQS_CHANGE_QUEUE",
        "name": "SQS_CHANGE_QUEUE"
      },
      {
        "valueFrom": "/${var.config_path}/DB_NAME",
        "name": "DB_NAME"
      },
      {
        "valueFrom": "/${var.config_path}/OAUTH_CLIENT_ID",
        "name": "OAUTH_CLIENT_ID"
      },
      {
        "valueFrom": "/${var.config_path}/OAUTH_CLIENT_SECRET",
        "name": "OAUTH_CLIENT_SECRET"
      },
      {
        "valueFrom": "/${var.config_path}/OAUTH_SESSION_SIGN",
        "name": "OAUTH_SESSION_SIGN"
      },
      {
        "valueFrom": "/${var.config_path}/OAUTH_TOKEN_URL",
        "name": "OAUTH_TOKEN_URL"
      },
      {
        "valueFrom": "/${var.config_path}/SQS_ALARM_QUEUE",
        "name": "SQS_ALARM_QUEUE"
      },
      {
        "valueFrom": "/${var.config_path}/OAUTH_AUTH_URL",
        "name": "OAUTH_AUTH_URL"
      },
      {
        "valueFrom": "/${var.config_path}/SNS_ALARM_TOPIC_ARN",
        "name": "SNS_ALARM_TOPIC_ARN"
      },
      {
        "valueFrom": "/${var.config_path}/SQS_S3_QUEUE",
        "name": "SQS_S3_QUEUE"
      },
      {
        "valueFrom": "/${var.config_path}/OAUTH_REDIRECT_URL",
        "name": "OAUTH_REDIRECT_URL"
      },
      {
        "valueFrom": "/${var.config_path}/OAUTH_SCOPES",
        "name": "OAUTH_SCOPES"
      }
    ],
    "logConfiguration": {
      "logDriver": "awslogs",
      "options": {
        "awslogs-group": "/fargate/service/${var.app}-${var.environment}",
        "awslogs-region": "us-east-1",
        "awslogs-stream-prefix": "ecs"
      }
    }
  }
]
DEFINITION

}

resource "aws_ecs_service" "app" {
  name            = "${var.app}-${var.environment}"
  cluster         = aws_ecs_cluster.app.id
  launch_type     = "FARGATE"
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = var.replicas

  network_configuration {
    security_groups = [aws_security_group.nsg_task.id]
    subnets         = split(",", var.private_subnets)
  }

  load_balancer {
    target_group_arn = aws_alb_target_group.main.id
    container_name   = var.container_name
    container_port   = var.container_port
  }

  # workaround for https://github.com/hashicorp/terraform/issues/12634
  depends_on = [aws_alb_listener.http]

  # [after initial apply] don't override changes made to task_definition
  # from outside of terrraform (i.e.; fargate cli)
  # lifecycle {
  #   ignore_changes = [task_definition]
  # }
}

# https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_execution_IAM_role.html
resource "aws_iam_role" "ecsTaskExecutionRole" {
  name               = "${var.app}-${var.environment}-ecs"
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
}

data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy_attachment" "ecsTaskExecutionRole_policy" {
  role       = aws_iam_role.ecsTaskExecutionRole.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "ecsTaskExecutionRoleConfig_policy" {
  role       = aws_iam_role.ecsTaskExecutionRole.name
  policy_arn = aws_iam_policy.config-policy.arn
}

resource "aws_iam_policy" "config-policy" {
  name   = "${var.app}-${var.environment}"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["ssm:GetParameters"],
      "Resource": "arn:aws:ssm:us-east-1:${data.aws_caller_identity.current.account_id}:parameter/${var.config_path}/*"
    }
  ]
}
EOF

}

resource "aws_cloudwatch_log_group" "logs" {
  name              = "/fargate/service/${var.app}-${var.environment}"
  retention_in_days = var.ecs_cloudwatch_log_retention_in_days
  tags              = var.tags
}

output "ecs_task_execution_role_arn" {
  value = aws_iam_role.ecsTaskExecutionRole.arn
}

output "app_role_arn" {
  value = aws_iam_role.app_role.arn
}

