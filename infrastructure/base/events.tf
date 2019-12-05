resource "aws_cloudwatch_event_permission" "linked_accounts" {
  count = length(var.linked_account_ids)

  principal  = var.linked_account_ids[count.index]
  statement_id = "LinkedAccountEvents-${count.index}"
}