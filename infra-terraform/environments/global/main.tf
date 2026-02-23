locals {
  tags = {
    Project     = "order-bot"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

module "ecr" {
  source = "../../modules/ecr"

  repository_names = [
    "order-bot-svc",
    "order-bot-mgmt-svc",
    "order-bot-frontend"
  ]
  tags = locals.tags
}

module "frontend" {
  source = "../../modules/s3"

  name_prefix         = "order-bot-${var.environment}"
  bucket_name         = var.frontend_bucket_name
  frontend_domain     = var.frontend_domain
  hosted_zone_id      = var.hosted_zone_id
  acm_certificate_arn = var.cloudfront_certificate_arn
  tags                = locals.tags
}
