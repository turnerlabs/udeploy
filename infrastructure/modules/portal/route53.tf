data "aws_route53_zone" "route_zone" {
  name         = var.zone_name
}

resource "aws_route53_record" "env" {
  zone_id = data.aws_route53_zone.route_zone.zone_id
  name    = var.record_name
  type    = "A"

  alias {
    name                   = aws_alb.main.dns_name
    zone_id                = aws_alb.main.zone_id
    evaluate_target_health = true
  }
}

data "aws_acm_certificate" "cert" {
  domain   = var.zone_name
}

output "record_name" {
  value = aws_route53_record.env.name
}

