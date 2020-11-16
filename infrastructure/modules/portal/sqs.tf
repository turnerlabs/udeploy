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

  policy = data.aws_iam_policy_document.alarm_queue_policy.json
}

data "aws_iam_policy_document" "alarm_queue_policy" {
  statement {
    sid = "allow portal account to publish alarms"
    effect = "Allow"

    actions = ["sqs:SendMessage"]

    resources = [aws_sqs_queue.alarm_queue.arn]

    principals {
      type = "AWS"
      identifiers = ["*"]
    }

    condition {
      test     = "ArnEquals"
      variable = "aws:SourceArn"

      values = [aws_sns_topic.alarms.arn]
    }

  }
}

resource "aws_sqs_queue" "s3_queue" {
  name = "${var.app}-${var.environment}-s3-queue"

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": { "AWS": ${jsonencode(concat(data.template_file.linked_account_ids.*.rendered, ["${data.aws_caller_identity.current.account_id}"]))} },
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

//render dynamic list of linked account ids
data "template_file" "linked_account_ids" {
  count    = length(var.linked_account_ids)
  template = "$${account}"

  vars = {
    account = var.linked_account_ids[count.index]
  }
}