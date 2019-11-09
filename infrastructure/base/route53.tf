resource "aws_route53_zone" "route_zone" {
  name = var.domain

  tags = var.tags
}

resource "aws_route53_record" "validation" {
  name    = lookup(aws_acm_certificate.cert.domain_validation_options[0], "resource_record_name")
  type    = lookup(aws_acm_certificate.cert.domain_validation_options[0], "resource_record_type")
  zone_id = aws_route53_zone.route_zone.zone_id
  records = [lookup(aws_acm_certificate.cert.domain_validation_options[0], "resource_record_value")]
  ttl     = 60
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

resource "aws_acm_certificate_validation" "main" {
  certificate_arn         = aws_acm_certificate.cert.arn
  validation_record_fqdns = [aws_route53_record.validation.fqdn]
}

output "zone_id" {
  value = aws_route53_zone.route_zone.zone_id
}

