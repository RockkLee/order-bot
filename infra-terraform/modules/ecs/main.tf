locals {
  containers = {
    orderbot = {
      name = "order-bot-svc"
      port = var.order_bot_port
    },
    orderbot_mgmt = {
      name = "order-bot-mgmt-svc"
      port = var.order_bot_mgmt_port
    }
  }
}

data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

resource "aws_ecs_cluster" "this" {
  name = "${var.name_prefix}-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = var.tags
}

# Allow the ECS Tasks service to assume (use) this IAM role.
data "aws_iam_policy_document" "ecs_task_assume" {
  statement {
    actions = ["sts:AssumeRole"] # the "sts:AssumeRole" action allows someone to assume (use) this role

    # This defines who is allowed to assume the role
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

# For aws_ecs_task_definition.execution_role_arn
# execution_role_arn: used by ECS/Fargate agent (pull image from ECR, send logs to CloudWatch, fetch secrets refs, etc.)
resource "aws_iam_role" "task_execution" {
  name               = "${var.name_prefix}-ecs-task-exec"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = var.tags
}

resource "aws_iam_role_policy_attachment" "task_exec_default" {
  role       = aws_iam_role.task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# For aws_ecs_task_definition.task_role_arn
# task_role_arn: used by your application code inside the container to call AWS APIs (S3, SQS, DynamoDB, Secrets Manager, etc.)
# (No policy is attached to it for now)
resource "aws_iam_role" "task_role" {
  name               = "${var.name_prefix}-ecs-task-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = var.tags
}

resource "aws_ecs_task_definition" "orderbot" {
  family                   = "${var.name_prefix}-order-bot-svc"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = tostring(var.order_bot_task_cpu)
  memory                   = tostring(var.order_bot_task_memory)
  execution_role_arn       = aws_iam_role.task_execution.arn
  task_role_arn            = aws_iam_role.task_role.arn

  container_definitions = jsonencode([
    {
      name      = local.containers.orderbot.name
      image     = var.order_bot_image
      cpu       = var.order_bot_task_cpu
      memory    = var.order_bot_task_memory
      essential = true
      portMappings = [{
        containerPort = local.containers.orderbot.port
        hostPort      = local.containers.orderbot.port
        protocol      = "tcp"
      }]
      environment = [
        for k, v in var.order_bot_environment : { name = k, value = v }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.orderbot.name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])

  tags = var.tags
}

resource "aws_ecs_task_definition" "orderbot_mgmt" {
  family                   = "${var.name_prefix}-order-bot-mgmt-svc"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = tostring(var.order_bot_mgmt_task_cpu)
  memory                   = tostring(var.order_bot_mgmt_task_memory)
  execution_role_arn       = aws_iam_role.task_execution.arn
  task_role_arn            = aws_iam_role.task_role.arn

  container_definitions = jsonencode([
    {
      name      = local.containers.orderbot_mgmt.name
      image     = var.order_bot_mgmt_image
      cpu       = var.order_bot_mgmt_task_cpu
      memory    = var.order_bot_mgmt_task_memory
      essential = true
      portMappings = [{
        containerPort = local.containers.orderbot_mgmt.port
        hostPort      = local.containers.orderbot_mgmt.port
        protocol      = "tcp"
      }]
      environment = [
        for k, v in var.order_bot_mgmt_environment : { name = k, value = v }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.orderbot_mgmt.name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])

  tags = var.tags
}

resource "aws_ecs_service" "orderbot" {
  name            = local.containers.orderbot.name
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.orderbot.arn
  desired_count   = var.order_bot_desired_count
  launch_type     = "FARGATE"

  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200

  lifecycle {
    precondition {
      condition     = var.enable_alb ? var.order_bot_target_group_arn != null : true
      error_message = "order_bot_target_group_arn is required when enable_alb is true."
    }
  }

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.app_security_group_id]
    assign_public_ip = false
  }

  # In Terraform, a dynamic block repeats once per item in for_each
  dynamic "load_balancer" {
    for_each = var.enable_alb ? [1] : []  # [] means “iterate once if enabled, or zero times if disabled.”
    content {
      target_group_arn = var.order_bot_target_group_arn
      container_name   = local.containers.orderbot.name
      container_port   = local.containers.orderbot.port
    }
  }

  # Bind the service discovery config to register an internal domain name
  # to allow other ECS instances to call this instance
  service_registries {
    registry_arn = aws_service_discovery_service.orderbot.arn

    # MUST match a container name in your task definition
    container_name = local.containers.orderbot.name
  }

  tags = var.tags
}

resource "aws_ecs_service" "orderbot_mgmt" {
  name            = local.containers.orderbot_mgmt.name
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.orderbot_mgmt.arn
  desired_count   = var.order_bot_mgmt_desired_count
  launch_type     = "FARGATE"

  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200

  lifecycle {
    precondition {
      condition     = var.enable_alb ? var.order_bot_mgmt_target_group_arn != null : true
      error_message = "order_bot_mgmt_target_group_arn is required when enable_alb is true."
    }
  }

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.app_security_group_id]
    assign_public_ip = false
  }

  # In Terraform, a dynamic block repeats once per item in for_each
  dynamic "load_balancer" {
    for_each = var.enable_alb ? [1] : []  # [] means “iterate once if enabled, or zero times if disabled.”
    content {
      target_group_arn = var.order_bot_mgmt_target_group_arn
      container_name   = local.containers.orderbot_mgmt.name
      container_port   = local.containers.orderbot_mgmt.port
    }
  }

  # Bind the service discovery config to register an internal domain name
  # to allow other ECS instances to call this instance
  service_registries {
    registry_arn = aws_service_discovery_service.orderbot_mgmt.arn

    # MUST match a container name in your task definition
    container_name = local.containers.orderbot_mgmt.name
  }

  tags = var.tags
}

resource "aws_appautoscaling_target" "orderbot" {
  max_capacity       = 1
  min_capacity       = 0
  resource_id        = "service/${aws_ecs_cluster.this.name}/${aws_ecs_service.orderbot.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_target" "orderbot_mgmt" {
  max_capacity       = 1
  min_capacity       = 0
  resource_id        = "service/${aws_ecs_cluster.this.name}/${aws_ecs_service.orderbot_mgmt.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}
