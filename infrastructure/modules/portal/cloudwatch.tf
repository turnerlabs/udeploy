resource "aws_cloudwatch_event_rule" "log_notification" {
  name          = "${var.app}-${var.environment}-log-notification"
  description   = "Notifies the queue when task definitions are registered."
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

resource "aws_cloudwatch_event_target" "log_sqs" {
  rule = aws_cloudwatch_event_rule.log_notification.name
  sqs_target {
    message_group_id = "event_group"
  }
  target_id = "SendToSQS"
  arn       = aws_sqs_queue.notification_queue.arn
}

resource "aws_cloudwatch_event_rule" "change_notification" {
  name        = "${var.app}-${var.environment}-change-notification"
  description = "Notifies the queue when tasks are changing status in ECS."
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

resource "aws_cloudwatch_event_target" "sqs" {
  rule = aws_cloudwatch_event_rule.change_notification.name
  sqs_target {
    message_group_id = "event_group"
  }
  target_id = "SendToSQS"
  arn       = aws_sqs_queue.notification_queue.arn
}

resource "aws_cloudwatch_event_rule" "lambda_notification" {
  name        = "${var.app}-${var.environment}-lambda-notification"
  description = "Notifies the queue when labmda events happen."

  event_pattern = <<PATTERN
{
  "source": [
    "aws.lambda"
  ]
}
PATTERN

}

resource "aws_cloudwatch_event_target" "lambda_sqs" {
  rule = aws_cloudwatch_event_rule.lambda_notification.name
  sqs_target {
    message_group_id = "event_group"
  }
  target_id = "SendToSQS"
  arn       = aws_sqs_queue.notification_queue.arn
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

resource "aws_sns_topic_policy" "alarms" {
  arn = aws_sns_topic.alarms.arn

  policy = data.aws_iam_policy_document.alarms.json
}

data "aws_iam_policy_document" "alarms" {
  policy_id = "${var.app}-${var.environment}-alarms"

  statement {
    sid = "allow linked accounts to publish to alarm topic"
    effect = "Allow"
    actions = ["SNS:Publish"]

    principals {
      type        = "AWS"
      identifiers = ["*"]
    }

    resources = [
      aws_sns_topic.alarms.arn,
    ]

    condition {
      test     = "ArnLike"
      variable = "aws:SourceArn"

      values = [
        for account_id in var.linked_account_ids:
          "arn:aws:cloudwatch:${var.region}:${account_id}:alarm:*"
      ]
    }     
  }
}

resource "aws_sns_topic_subscription" "lambda_alerts_sqs_target" {
  topic_arn = aws_sns_topic.alarms.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.alarm_queue.arn
}