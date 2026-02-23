locals {
  tags = {
    Project     = "order-bot"
    Environment = "prod"
    ManagedBy   = "terraform"
  }
}

module "security_group" {
  source = "../../modules/security_group"

  name_prefix       = "order-bot-prod"
  vpc_id            = var.vpc_id
  app_port          = var.order_bot_port
  alb_ingress_cidrs = var.alb_ingress_cidrs
  tags              = locals.tags
}

module "alb" {
  source = "../../modules/alb"

  name_prefix          = "order-bot-prod"
  vpc_id               = var.vpc_id
  public_subnet_ids    = var.public_subnet_ids
  alb_security_group_id  = module.security_group.alb_security_group_id
  acm_certificate_arn  = var.alb_certificate_arn
  mgmt_host            = var.mgmt_domain
  chat_host            = var.chat_domain
  order_bot_port       = var.order_bot_port
  order_bot_mgmt_port  = var.order_bot_mgmt_port
  tags                 = locals.tags
}

resource "aws_route53_record" "mgmt_alias" {
  zone_id = var.hosted_zone_id
  name    = var.mgmt_domain
  type    = "A"

  alias {
    name                   = module.alb.alb_dns_name
    zone_id                = module.alb.alb_zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "chat_alias" {
  zone_id = var.hosted_zone_id
  name    = var.chat_domain
  type    = "A"

  alias {
    name                   = module.alb.alb_dns_name
    zone_id                = module.alb.alb_zone_id
    evaluate_target_health = true
  }
}

module "ecs" {
  source = "../../modules/ecs"

  name_prefix                    = "order-bot-prod"
  private_subnet_ids             = var.private_subnet_ids
  app_security_group_id          = module.security_group.app_security_group_id
  order_bot_target_group_arn     = module.alb.order_bot_target_group_arn
  order_bot_mgmt_target_group_arn = module.alb.order_bot_mgmt_target_group_arn

  order_bot_image      = var.order_bot_image
  order_bot_mgmt_image = var.order_bot_mgmt_image

  order_bot_port      = var.order_bot_port
  order_bot_mgmt_port = var.order_bot_mgmt_port

  order_bot_environment      = var.order_bot_environment
  order_bot_mgmt_environment = var.order_bot_mgmt_environment

  order_bot_desired_count      = var.default_desired_count
  order_bot_mgmt_desired_count = var.default_desired_count

  tags = locals.tags
}

module "schedule" {
  source = "../../modules/schedule"

  resources = {
    "order-bot-svc"      = module.ecs.order_bot_autoscaling_resource_id
    "order-bot-mgmt-svc" = module.ecs.order_bot_mgmt_autoscaling_resource_id
  }

  timezone        = var.autoscaling_timezone
  scale_up_cron   = var.scale_up_cron
  scale_down_cron = var.scale_down_cron
}
