resource "aws_cloudwatch_event_rule" "log_notification" {
  name          = "${var.app}-${var.environment}-log-notification"
  description   = "Watch CloudTrail logs for task definition registrations."
  is_enabled    = true
  event_pattern = <<PATTERN
{
  "source": [
    "aws.ecs"
  ],
  "detail-type": [
    "AWS API Call via CloudTrail"
  ],
  "detail": {
    "eventSource": [
      "ecs.amazonaws.com"
    ],
    "eventName": [
      "RegisterTaskDefinition"
    ]
  }
}
PATTERN

}

resource "aws_cloudwatch_event_target" "log-event" {
  rule      = aws_cloudwatch_event_rule.log_notification.name
  target_id = aws_cloudwatch_event_rule.log_notification.name
  arn       = "arn:aws:events:${var.region}:${var.portal_account_id}:event-bus/default"
  role_arn  = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${var.app}-${var.environment}"
}

resource "aws_cloudwatch_event_rule" "task_notification" {
  name        = "${var.app}-${var.environment}-task-notification"
  description = "Watch for task status changes in ECS."
  is_enabled  = true

  event_pattern = <<PATTERN
{
  "source": [
    "aws.ecs"
  ],
  "detail-type": [
    "ECS Task State Change"
  ]
}
PATTERN

}

resource "aws_cloudwatch_event_target" "task-event" {
  rule      = aws_cloudwatch_event_rule.task_notification.name
  target_id = aws_cloudwatch_event_rule.task_notification.name
  arn       = "arn:aws:events:${var.region}:${var.portal_account_id}:event-bus/default"
  role_arn  = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${var.app}-${var.environment}"
}

resource "aws_cloudwatch_event_rule" "lambda_notification" {
  name        = "${var.app}-${var.environment}-lambda-notification"
  description = "Watch for lambda events."

  event_pattern = <<PATTERN
{
  "source": [
    "aws.lambda"
  ]
}
PATTERN

}

resource "aws_cloudwatch_event_target" "lambda-event" {
  rule      = aws_cloudwatch_event_rule.lambda_notification.name
  target_id = aws_cloudwatch_event_rule.lambda_notification.name
  arn       = "arn:aws:events:${var.region}:${var.portal_account_id}:event-bus/default"
  role_arn  = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/${var.app}-${var.environment}"
}

resource "aws_sns_topic" "alarms" {
  name = "${var.app}-${var.environment}-alarms"

  delivery_policy = <<EOF
{
  "http": {
    "defaultHealthyRetryPolicy": {
      "minDelayTarget": 20,
      "maxDelayTarget": 20,
      "numRetries": 3,
      "numMaxDelayRetries": 0,
      "numNoDelayRetries": 0,
      "numMinDelayRetries": 0,
      "backoffFunction": "linear"
    },
    "disableSubscriptionOverrides": false,
    "defaultThrottlePolicy": {
      "maxReceivesPerSecond": 1
    }
  }
}
EOF

}