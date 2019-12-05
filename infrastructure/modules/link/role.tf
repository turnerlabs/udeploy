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
      "events:PutEvents"
    ]

    resources = [
      "arn:aws:events:${var.region}:${var.portal_account_id}:event-bus/default",
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
      type = "AWS"
 
      identifiers = [
        "arn:aws:iam::${var.portal_account_id}:role/${var.app}-${var.environment}",
      ]
    }
  }

  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type = "Service"
 
      identifiers = [
        "events.amazonaws.com",
      ]
    } 
  }
}

output "account_id" {
  value = "${data.aws_caller_identity.current.account_id}"
}

output "role_arn" {
  value = "${aws_iam_role.app_role.arn}"
}