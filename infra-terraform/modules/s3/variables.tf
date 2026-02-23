variable "name_prefix" { type = string }
variable "bucket_name" { type = string }
variable "frontend_domain" { type = string }
variable "hosted_zone_id" { type = string }
variable "acm_certificate_arn" { type = string }
variable "tags" { type = map(string) default = {} }
