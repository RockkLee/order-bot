output "ecr_repository_urls" {
  value = module.ecr.repository_urls
}

output "frontend_bucket_name" {
  value = module.frontend.bucket_name
}

output "frontend_cloudfront_domain" {
  value = module.frontend.cloudfront_domain_name
}
