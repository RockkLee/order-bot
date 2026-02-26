resource "aws_s3_bucket" "frontend" {
  bucket = var.bucket_name
  tags   = merge(var.tags, { Name = var.bucket_name })
}

resource "aws_s3_bucket_public_access_block" "frontend" {
  bucket                  = aws_s3_bucket.frontend.id
  block_public_acls       = true  # Block public ACLs
  block_public_policy     = true  # Block public policies
  ignore_public_acls      = true  # Even if public ACLs already exist, S3 ignores them for access evaluation.
  # if TRUE: restricts access to this bucket to only Amazon Web Services service principals and authorized users.
  restrict_public_buckets = true
}

# OAC: Origin Access Control: a new feature that is similar to OAI
# Cloudfront: CDN / Cloudfront origin: where the file originally comes from
resource "aws_cloudfront_origin_access_control" "this" {
  name                              = "${var.name_prefix}-oac"
  description                       = "OAC for frontend bucket"
  origin_access_control_origin_type = "s3"

  # CloudFront sends the request
  # It includes a calculated signature (SigV4)
  # S3 verifies it using AWS auth rules
  # S3 can trust: “Yep, this request is legit”
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

data "aws_cloudfront_cache_policy" "caching_optimized" {
  name = "Managed-CachingOptimized"
}

resource "aws_cloudfront_distribution" "this" {
  enabled             = true
  default_root_object = "index.html"
  aliases             = [var.frontend_domain]

  origin {
    domain_name              = aws_s3_bucket.frontend.bucket_regional_domain_name
    origin_id                = "frontend-s3"
    origin_access_control_id = aws_cloudfront_origin_access_control.this.id
  }

  default_cache_behavior {
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = "frontend-s3"
    viewer_protocol_policy = "redirect-to-https"
    cache_policy_id        = data.aws_cloudfront_cache_policy.caching_optimized.id
    # Replace `forwarded_values` (deprecated) with `cache_policy_id`
    # forwarded_values {
    #   query_string = false
    #   cookies {
    #     forward = "none"
    #   }
    # }
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  # The SSL config for this distribution
  viewer_certificate {
    acm_certificate_arn      = var.acm_certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }

  # custom_error_response {
  #   error_code            = 403
  #   response_code         = 200
  #   response_page_path    = "/index.html"
  #   error_caching_min_ttl = 0
  # }

  # custom_error_response {
  #   error_code            = 404
  #   response_code         = 200
  #   response_page_path    = "/index.html"
  #   error_caching_min_ttl = 0
  # }

  tags = var.tags
}

resource "aws_s3_bucket_policy" "frontend" {
  bucket = aws_s3_bucket.frontend.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        # Sid (statement ID):
        Sid       = "AllowCloudFrontServicePrincipalReadOnly"
        Effect    = "Allow"
        Principal = { Service = "cloudfront.amazonaws.com" }
        Action    = ["s3:GetObject"]
        Resource  = ["${aws_s3_bucket.frontend.arn}/*"]
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.this.arn
          }
        }
      }
    ]
  })
}

resource "aws_route53_record" "frontend_alias" {
  zone_id = var.hosted_zone_id
  name    = var.frontend_domain
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.this.domain_name
    zone_id                = aws_cloudfront_distribution.this.hosted_zone_id
    evaluate_target_health = false
  }
}
