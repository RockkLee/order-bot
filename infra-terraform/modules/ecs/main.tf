data "aws_region" "current" {}

data "aws_caller_identity" "current" {}

resource "aws_cloudwatch_log_group" "svc" {
  name              = "/ecs/${var.name_prefix}/order-bot-svc"
  retention_in_days = 30
  tags              = var.tags
}

resource "aws_cloudwatch_log_group" "mgmt" {
  name              = "/ecs/${var.name_prefix}/order-bot-mgmt-svc"
  retention_in_days = 30
  tags              = var.tags
}

resource "aws_ecs_cluster" "this" {
  name = "${var.name_prefix}-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = var.tags
}

data "aws_iam_policy_document" "ecs_task_assume" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "task_execution" {
  name               = "${var.name_prefix}-ecs-task-exec"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = var.tags
}

resource "aws_iam_role_policy_attachment" "task_exec_default" {
  role       = aws_iam_role.task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role" "task_role" {
  name               = "${var.name_prefix}-ecs-task-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_assume.json
  tags               = var.tags
}

resource "aws_ecs_task_definition" "svc" {
  family                   = "${var.name_prefix}-order-bot-svc"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = tostring(var.order_bot_task_cpu)
  memory                   = tostring(var.order_bot_task_memory)
  execution_role_arn       = aws_iam_role.task_execution.arn
  task_role_arn            = aws_iam_role.task_role.arn

  container_definitions = jsonencode([
    {
      name      = "order-bot-svc"
      image     = var.order_bot_image
      cpu       = var.order_bot_task_cpu
      memory    = var.order_bot_task_memory
      essential = true
      portMappings = [{
        containerPort = var.order_bot_port
        hostPort      = var.order_bot_port
        protocol      = "tcp"
      }]
      environment = [
        for k, v in var.order_bot_environment : { name = k, value = v }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.svc.name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])

  tags = var.tags
}

resource "aws_ecs_task_definition" "mgmt" {
  family                   = "${var.name_prefix}-order-bot-mgmt-svc"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = tostring(var.order_bot_mgmt_task_cpu)
  memory                   = tostring(var.order_bot_mgmt_task_memory)
  execution_role_arn       = aws_iam_role.task_execution.arn
  task_role_arn            = aws_iam_role.task_role.arn

  container_definitions = jsonencode([
    {
      name      = "order-bot-mgmt-svc"
      image     = var.order_bot_mgmt_image
      cpu       = var.order_bot_mgmt_task_cpu
      memory    = var.order_bot_mgmt_task_memory
      essential = true
      portMappings = [{
        containerPort = var.order_bot_mgmt_port
        hostPort      = var.order_bot_mgmt_port
        protocol      = "tcp"
      }]
      environment = [
        for k, v in var.order_bot_mgmt_environment : { name = k, value = v }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = aws_cloudwatch_log_group.mgmt.name
          awslogs-region        = data.aws_region.current.name
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])

  tags = var.tags
}

resource "aws_ecs_service" "svc" {
  name            = "order-bot-svc"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.svc.arn
  desired_count   = var.order_bot_desired_count
  launch_type     = "FARGATE"

  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.app_security_group_id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = var.order_bot_target_group_arn
    container_name   = "order-bot-svc"
    container_port   = var.order_bot_port
  }

  tags = var.tags
}

resource "aws_ecs_service" "mgmt" {
  name            = "order-bot-mgmt-svc"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.mgmt.arn
  desired_count   = var.order_bot_mgmt_desired_count
  launch_type     = "FARGATE"

  deployment_minimum_healthy_percent = 100
  deployment_maximum_percent         = 200

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = [var.app_security_group_id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = var.order_bot_mgmt_target_group_arn
    container_name   = "order-bot-mgmt-svc"
    container_port   = var.order_bot_mgmt_port
  }

  tags = var.tags
}

resource "aws_appautoscaling_target" "svc" {
  max_capacity       = 1
  min_capacity       = 0
  resource_id        = "service/${aws_ecs_cluster.this.name}/${aws_ecs_service.svc.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_target" "mgmt" {
  max_capacity       = 1
  min_capacity       = 0
  resource_id        = "service/${aws_ecs_cluster.this.name}/${aws_ecs_service.mgmt.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}
