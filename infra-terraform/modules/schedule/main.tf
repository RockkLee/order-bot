resource "aws_appautoscaling_scheduled_action" "scale_up" {
  for_each = var.resources

  name               = "${each.key}-scale-up"
  service_namespace  = "ecs"
  scalable_dimension = "ecs:service:DesiredCount"
  resource_id        = each.value
  schedule           = var.scale_up_cron
  timezone           = var.timezone

  scalable_target_action {
    min_capacity = 1
    max_capacity = 1
  }
}

resource "aws_appautoscaling_scheduled_action" "scale_down" {
  for_each = var.resources

  name               = "${each.key}-scale-down"
  service_namespace  = "ecs"
  scalable_dimension = "ecs:service:DesiredCount"
  resource_id        = each.value
  schedule           = var.scale_down_cron
  timezone           = var.timezone

  scalable_target_action {
    min_capacity = 0
    max_capacity = 1
  }
}
