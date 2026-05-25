locals {
  endpoint_policy = jsonencode({
    Statement = [{
      Action    = "*"
      Effect    = "Allow"
      Principal = "*"
      Resource  = "*"
    }]
  })
}

resource "aws_vpc_endpoint" "ecr_api" {
  ip_address_type     = "ipv4"
  policy              = local.endpoint_policy
  private_dns_enabled = true
  region              = var.region
  security_group_ids  = var.security_group_ids
  service_name        = "com.amazonaws.${var.region}.ecr.api"
  service_region      = var.region
  subnet_ids          = var.subnet_ids
  tags                = merge(var.tags, { Name = "endpoint-com.amazonaws.${var.region}.ecr.api" })
  vpc_endpoint_type   = "Interface"
  vpc_id              = var.vpc_id
}

resource "aws_vpc_endpoint" "ecr_dkr" {
  ip_address_type     = "ipv4"
  policy              = local.endpoint_policy
  private_dns_enabled = true
  region              = var.region
  security_group_ids  = var.security_group_ids
  service_name        = "com.amazonaws.${var.region}.ecr.dkr"
  service_region      = var.region
  subnet_ids          = var.subnet_ids
  tags                = merge(var.tags, { Name = "endpoint-com.amazonaws.${var.region}.ecr.dkr" })
  vpc_endpoint_type   = "Interface"
  vpc_id              = var.vpc_id
}

resource "aws_vpc_endpoint" "logs" {
  ip_address_type     = "ipv4"
  policy              = local.endpoint_policy
  private_dns_enabled = true
  region              = var.region
  security_group_ids  = var.security_group_ids
  service_name        = "com.amazonaws.${var.region}.logs"
  service_region      = var.region
  subnet_ids          = var.subnet_ids
  tags                = merge(var.tags, { Name = "endpoint-com.amazonaws.${var.region}.logs" })
  vpc_endpoint_type   = "Interface"
  vpc_id              = var.vpc_id
}
