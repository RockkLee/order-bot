output "cluster_name" { value = aws_ecs_cluster.this.name }
output "order_bot_autoscaling_resource_id" { value = aws_appautoscaling_target.svc.resource_id }
output "order_bot_mgmt_autoscaling_resource_id" { value = aws_appautoscaling_target.mgmt.resource_id }
