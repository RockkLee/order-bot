variable "aws_region" { type = string }
variable "vpc_id" { type = string }
variable "public_subnet_ids" { type = list(string) }
variable "private_subnet_ids" { type = list(string) }
variable "hosted_zone_id" { type = string }

variable "alb_certificate_arn" { type = string }
variable "orderbot_domain" { type = string }
variable "orderbot_mgmt_domain" { type = string }

variable "order_bot_image" { type = string }
variable "order_bot_mgmt_image" { type = string }

variable "order_bot_port" {
  type    = number
  default = 8000
}
variable "order_bot_mgmt_port" {
  type    = number
  default = 8080
}

variable "default_desired_count" {
  type    = number
  default = 1
}

# cidrs: Classless Inter‑Domain Routing
# A fancy name that refers to IPv4 or IPv6
variable "alb_ingress_cidrs" {
  type    = list(string)
  default = ["0.0.0.0/0"]
}

variable "autoscaling_timezone" {
  type    = string
  default = "Asia/Tokyo"
}

variable "scale_up_cron" {
  type    = string
  default = "cron(0 10 * * ? *)"
}

variable "scale_down_cron" {
  type    = string
  default = "cron(0 18 * * ? *)"
}

variable "order_bot_environment" {
  description = "Environment variables for order-bot-svc (FastAPI)"
  type        = map(string)
}

variable "order_bot_mgmt_environment" {
  description = "Environment variables for order-bot-mgmt-svc (Go)"
  type        = map(string)
}
