resource "aws_sqs_queue" "notification_queue" {
  name                        = "${var.app}-${var.environment}-notification-queue.fifo"
  fifo_queue                  = true
  content_based_deduplication = true
}

resource "aws_sqs_queue_policy" "notification_queue_policy" {
  queue_url = aws_sqs_queue.notification_queue.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "sqspolicy",
  "Statement": [
    {
      "Sid": "First",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.notification_queue.arn}",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_cloudwatch_event_rule.change_notification.arn}"
        }
      }
    },
    {
      "Sid": "First",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.notification_queue.arn}",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_cloudwatch_event_rule.log_notification.arn}"
        }
      }
    },
    {
      "Sid": "First",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.notification_queue.arn}",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_cloudwatch_event_rule.lambda_notification.arn}"
        }
      }
    }
  ]
}
POLICY

}

resource "aws_sqs_queue" "alarm_queue" {
  name = "${var.app}-${var.environment}-alarm-queue"
}

resource "aws_sqs_queue_policy" "alarm_queue_policy" {
  queue_url = aws_sqs_queue.alarm_queue.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "sqspolicy",
  "Statement": [
    {
      "Sid": "First",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.alarm_queue.arn}",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_sns_topic.alarms.arn}"
        }
      }
    }
  ]
}
POLICY

}

resource "aws_sqs_queue" "s3_queue" {
  name = "${var.app}-${var.environment}-s3-queue"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "arn:aws:sqs:*:*:${var.app}-${var.environment}-s3-queue",
      "Condition": {
        "ArnEquals": { "aws:SourceArn": "arn:aws:s3:*:*:*" }
      }
    }
  ]
}
POLICY

}

# SQS queue that is watched for s3 deployment 
# changes when updating the portal ui.
output "s3_change_queue" {
  value = aws_sqs_queue.s3_queue.name
}

