variable "name_prefix" {
  type = string
}

variable "vpc_id" {
  type = string
}

variable "order_bot_port" {
  description = "Port exposed by ECS tasks"
  type        = number
}

variable "order_bot_mgmt_port" {
  description = "Port exposed by ECS tasks"
  type        = number
}

variable "alb_ingress_cidrs" {
  type    = list(string)
  default = ["0.0.0.0/0"]
}

variable "tags" {
  type    = map(string)
  default = {}
}
