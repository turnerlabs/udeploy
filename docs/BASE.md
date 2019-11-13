# Base Infrastructure
  
- [infrastructure/base/route53.tf](/infrastructure/base/route53.tf) (_optional_)
    Creates an AWS Route53 Hosted Zone, a certificate, and the DNS validation N records. `terraform apply` may take longer during validation. To complete validation, the Hosted Zone and DNS records must be publicly available on the internet. https://docs.aws.amazon.com/acm/latest/userguide/gs-acm-validate-dns.html
