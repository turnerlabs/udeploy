# The SAML role to use for adding users to the ECR policy
variable "saml_role" {
}

# creates an application role that the container/task runs as
resource "aws_iam_role" "app_role" {
  name               = "${var.app}-${var.environment}"
  assume_role_policy = data.aws_iam_policy_document.app_role_assume_role_policy.json
}

# assigns the app policy
resource "aws_iam_role_policy" "app_policy" {
  name   = "${var.app}-${var.environment}"
  role   = aws_iam_role.app_role.id
  policy = data.aws_iam_policy_document.app_policy.json
}

# TODO: fill out custom policy
data "aws_iam_policy_document" "app_policy" {
  statement {
    actions = [
      "ecs:DescribeClusters",
    ]

    resources = [
      aws_ecs_cluster.app.arn,
    ]
  }
  statement {
    actions = [
      "secretsmanager:GetSecretValue",
      "secretsmanager:CreateSecret",
      "secretsmanager:UpdateSecret",
      "secretsmanager:DeleteSecret",
    ]

    resources = [
      "arn:aws:secretsmanager:${var.region}:${data.aws_caller_identity.current.account_id}:secret:${var.app}/${var.environment}/*",
    ]
  }
  statement {
    actions = [
      "sqs:*",
    ]

    resources = [
      aws_sqs_queue.notification_queue.arn,
      aws_sqs_queue.alarm_queue.arn,
      aws_sqs_queue.s3_queue.arn,
    ]
  }
  statement {
    actions = [
      "ecs:DescribeServices",
      "ecs:DescribeTaskDefinition",
      "ecs:UpdateService",
      "ecs:RegisterTaskDefinition",
      "ecs:ListTaskDefinitions",
      "ecs:DescribeTasks",
      "ecs:ListTasks",
      "ecs:RunTask",
      "ecs:StopTask",
      "iam:PassRole",
      "events:DescribeRule",
      "events:ListTargetsByRule",
      "events:PutTargets",
      "application-autoscaling:DescribeScalableTargets",
      "ecr:ListImages",
      "lambda:GetAlias",
      "lambda:UpdateAlias",
      "lambda:DeleteFunction",
      "lambda:PublishVersion",
      "lambda:GetFunction",
      "lambda:GetFunctionConfiguration",
      "lambda:ListVersionsByFunction",
      "lambda:UpdateFunctionConfiguration",
      "lambda:UpdateFunctionCode",
      "lambda:InvokeFunction",
      "lambda:InvokeAsync",
      "sns:Publish",
      "cloudwatch:PutMetricAlarm",
      "cloudwatch:DescribeAlarms",
      "cloudwatch:DeleteAlarms",
      "s3:PutObject",
      "s3:GetObject",
      "s3:DeleteObject",
      "s3:ListBucket",
      "sts:AssumeRole",
      "cloudfront:CreateInvalidation",
      "cloudfront:GetInvalidation",
      "cloudfront:ListInvalidations",      
    ]

    resources = [
      "*",
    ]
  }
}

data "aws_caller_identity" "current" {
}

# allow role to be assumed by ecs and local saml users (for development)
data "aws_iam_policy_document" "app_role_assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }

    principals {
      type = "AWS"

      identifiers = [
        "arn:aws:sts::${data.aws_caller_identity.current.account_id}:assumed-role/${var.saml_role}/${var.saml_users[0]}",
      ]
    }
  }
}

