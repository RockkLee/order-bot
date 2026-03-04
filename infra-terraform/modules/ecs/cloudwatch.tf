resource "aws_cloudwatch_log_group" "orderbot" {
  name              = "/ecs/${var.name_prefix}/order-bot-svc"
  retention_in_days = 30
  tags              = var.tags
}

resource "aws_cloudwatch_log_group" "orderbot_mgmt" {
  name              = "/ecs/${var.name_prefix}/order-bot-mgmt-svc"
  retention_in_days = 30
  tags              = var.tags
}
