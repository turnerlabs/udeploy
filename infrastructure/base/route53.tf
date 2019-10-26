resource "aws_route53_zone" "route_zone" {
  name = var.domain

  tags = var.tags
}

resource "aws_route53_record" "route" {
  zone_id = aws_route53_zone.route_zone.zone_id
  name    = var.domain
  type    = "A"

  alias {
    name                   = var.alias_name
    zone_id                = var.alias_zone_id
    evaluate_target_health = false
  }
}

resource "aws_acm_certificate" "cert" {
  domain_name       = var.domain
  validation_method = "DNS"

  subject_alternative_names = [
    "*.${var.domain}",
  ]

  tags = var.tags

  lifecycle {
    create_before_destroy = true
  }
}

output "zone_id" {
  value = aws_route53_zone.route_zone.zone_id
}

