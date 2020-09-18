locals {
    # list of saml users for policies
    configUserIds = flatten([
        data.aws_caller_identity.current.account_id,
        "${aws_iam_role.app_role.unique_id}:*",
        formatlist(
        "%s:%s",
        data.aws_iam_role.saml_role_config.unique_id,
        var.saml_users,
        )
    ])

    # list of role users and saml users for policies
    configRoleIds = flatten([
        "${aws_iam_role.ecsTaskExecutionRole.unique_id}:*",
        "${aws_iam_role.app_role.unique_id}:*", 
    ])
}

# The users (email addresses) from the saml role to give access
# case sensitive
variable "saml_users" {
  type = list(string)
}

# get the saml user info so we can get the unique_id
data "aws_iam_role" "saml_role_config" {
  name = var.saml_role
}

resource "aws_kms_key" "config" {
  deletion_window_in_days = 7
  
  enable_key_rotation     = true
  
  tags = var.tags

  policy = data.template_file.config_policy.rendered
}

resource "aws_kms_alias" "config" {
  name          = "alias/${var.app}-${var.environment}"
  target_key_id = aws_kms_key.config.id
}

data "template_file" "config_policy" {
  template = <<EOF
  {
      "Version": "2012-10-17",
      "Statement": [
          {
              "Sid": "DenyWriteToAllExceptSAMLUsers",
              "Effect": "Deny",
              "Principal": {
                  "AWS": "*"
              },
              "Action": [
                  "kms:UpdateKeyDescription",
                  "kms:UpdateAlias",
                  "kms:UntagResource",
                  "kms:TagResource",
                  "kms:ScheduleKeyDeletion",
                  "kms:RevokeGrant",
                  "kms:RetireGrant",
                  "kms:ReEncryptTo",
                  "kms:ReEncryptFrom",
                  "kms:PutKeyPolicy",
                  "kms:ImportKeyMaterial",
                  "kms:GetParametersForImport",
                  "kms:GenerateRandom",
                  "kms:GenerateDataKeyWithoutPlaintext",
                  "kms:GenerateDataKey",
                  "kms:EnableKeyRotation",
                  "kms:EnableKey",
                  "kms:Encrypt",
                  "kms:DisableKeyRotation",
                  "kms:DisableKey",
                  "kms:DeleteImportedKeyMaterial",
                  "kms:DeleteAlias",
                  "kms:CreateKey",
                  "kms:CreateGrant",
                  "kms:CreateAlias",
                  "kms:CancelKeyDeletion"
              ],
              "Resource": "*",
              "Condition": {
                  "StringNotLike": {
                      "aws:userId": $${writePrincipals}
                  }
              }
          },
          {
              "Sid": "DenyReadToAllExceptRoleAndSAMLUsers",
              "Effect": "Deny",
              "Principal": {
                  "AWS": "*"
              },
              "Action": [
                  "kms:Decrypt"
              ],
              "Resource": "*",
              "Condition": {
                  "StringNotLike": {
                      "aws:userId": $${readPrincipals}
                  }
              }
          },
          {
              "Sid": "AllowWriteToSAMLUsers",
              "Effect": "Allow",
              "Principal": {
                  "AWS": "*"
              },
              "Action": [
                  "kms:UpdateKeyDescription",
                  "kms:UpdateAlias",
                  "kms:UntagResource",
                  "kms:TagResource",
                  "kms:ScheduleKeyDeletion",
                  "kms:RevokeGrant",
                  "kms:RetireGrant",
                  "kms:ReEncryptTo",
                  "kms:ReEncryptFrom",
                  "kms:PutKeyPolicy",
                  "kms:ImportKeyMaterial",
                  "kms:GetParametersForImport",
                  "kms:GenerateRandom",
                  "kms:GenerateDataKeyWithoutPlaintext",
                  "kms:GenerateDataKey",
                  "kms:EnableKeyRotation",
                  "kms:EnableKey",
                  "kms:Encrypt",
                  "kms:DisableKeyRotation",
                  "kms:DisableKey",
                  "kms:DeleteImportedKeyMaterial",
                  "kms:DeleteAlias",
                  "kms:CreateKey",
                  "kms:CreateGrant",
                  "kms:CreateAlias",
                  "kms:CancelKeyDeletion"
              ],
              "Resource": "*",
              "Condition": {
                  "StringLike": {
                      "aws:userId": $${writePrincipals}
                  }
              }
          },
          {
              "Sid": "AllowReadRoleAndSAMLUsers",
              "Effect": "Allow",
              "Principal": {
                  "AWS": "*"
              },
              "Action": [
                  "kms:Decrypt"
              ],
              "Resource": "*",
              "Condition": {
                  "StringLike": {
                      "aws:userId": $${readPrincipals}
                  }
              }
          },
          {
            "Sid": "Enable IAM User Permissions",
            "Effect": "Allow",
            "Principal": {
                "AWS": "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
            },
            "Action": [
                "kms:DescribeKey",
                "kms:GetKeyPolicy",
                "kms:GetKeyRotationStatus",
                "kms:List*"
            ],
            "Resource": "*"
          }
      ]
  }
  EOF

  vars = {
    writePrincipals = jsonencode(local.configUserIds)

    readPrincipals = jsonencode(
      concat(
        local.configUserIds,
        local.configRoleIds,
      ),
    )
  }
}

output "config_key_id" {
  value = aws_kms_key.config.key_id
}