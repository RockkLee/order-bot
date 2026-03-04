resource "aws_service_discovery_private_dns_namespace" "internal" {
  name        = "internal" # results in: *.internal
  description = "Service discovery for ECS services"
  vpc         = var.vpc_id
}

resource "aws_service_discovery_service" "orderbot" {
  name = "orderbot" # DNS name => svc2.internal

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.internal.id

    # For ECS tasks in awsvpc mode, use A records (task ENI IPs)
    dns_records {
      ttl  = 30
      type = "A"
    }

    routing_policy = "MULTIVALUE" # returns multiple task IPs
  }

  # Optional but recommended: only register healthy tasks
  # health_check_custom_config {
  #   failure_threshold = 1
  # }
}

resource "aws_service_discovery_service" "orderbot_mgmt" {
  name = "orderbotmgmt" # DNS name => svc2.internal

  dns_config {
    namespace_id = aws_service_discovery_private_dns_namespace.internal.id

    # For ECS tasks in awsvpc mode, use A records (task ENI IPs)
    dns_records {
      ttl  = 30
      type = "A"
    }

    routing_policy = "MULTIVALUE" # returns multiple task IPs
  }

  # Optional but recommended: only register healthy tasks
  # health_check_custom_config {
  #   failure_threshold = 1
  # }
}
