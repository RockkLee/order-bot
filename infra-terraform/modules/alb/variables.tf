variable "name_prefix" { type = string }
variable "vpc_id" { type = string }
variable "public_subnet_ids" { type = list(string) }
variable "alb_security_group_id" { type = string }
variable "acm_certificate_arn" { type = string }
variable "orderbot_mgmt_host" { type = string }
variable "orderbot_host" { type = string }
variable "order_bot_port" { type = number }
variable "order_bot_mgmt_port" { type = number }
variable "tags" { type = map(string) default = {} }
