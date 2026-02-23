variable "resources" {
  description = "Autoscaling resource IDs keyed by service name"
  type        = map(string)
}

variable "timezone" {
  type    = string
  default = "Asia/Tokyo"
}

variable "scale_up_cron" {
  description = "10:00"
  type        = string
  default     = "cron(0 10 * * ? *)"
}

variable "scale_down_cron" {
  description = "18:00"
  type        = string
  default     = "cron(0 18 * * ? *)"
}
