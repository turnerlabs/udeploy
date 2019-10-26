data "aws_route53_zone" "route_zone" {
  name         = var.zone_name
}

resource "aws_route53_record" "env" {
  zone_id = "${data.aws_route53_zone.route_zone.zone_id}"
  name    = var.record_name
  type    = "CNAME"
  ttl     = "60"

  records = ["${aws_alb.main.dns_name}"]
}

data "aws_acm_certificate" "cert" {
  domain   = var.zone_name
}

output "record_name" {
  value = aws_route53_record.env.name
}

