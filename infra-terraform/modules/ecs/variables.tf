variable "name_prefix" { type = string }
variable "private_subnet_ids" { type = list(string) }
variable "app_security_group_id" { type = string }
variable "order_bot_target_group_arn" { type = string }
variable "order_bot_mgmt_target_group_arn" { type = string }

variable "order_bot_image" { type = string }
variable "order_bot_mgmt_image" { type = string }

variable "order_bot_port" { type = number }
variable "order_bot_mgmt_port" { type = number }

variable "order_bot_task_cpu" { type = number default = 256 }
variable "order_bot_task_memory" { type = number default = 512 }
variable "order_bot_mgmt_task_cpu" { type = number default = 256 }
variable "order_bot_mgmt_task_memory" { type = number default = 512 }

variable "order_bot_desired_count" { type = number default = 1 }
variable "order_bot_mgmt_desired_count" { type = number default = 1 }

variable "order_bot_environment" {
  type    = map(string)
  default = {}
}

variable "order_bot_mgmt_environment" {
  type    = map(string)
  default = {}
}

variable "tags" {
  type    = map(string)
  default = {}
}
