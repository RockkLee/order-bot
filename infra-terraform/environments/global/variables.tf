variable "aws_region" { type = string }
variable "aws_profile" {
  type    = string
  default = null
}
variable "environment" {
  type    = string
  default = "global"
}
variable "hosted_zone_id" { type = string }
variable "frontend_domain" { type = string }
variable "frontend_bucket_name" { type = string }
variable "cloudfront_certificate_arn" { type = string }
