output "alb_dns_name" { value = aws_lb.this.dns_name }
output "alb_zone_id" { value = aws_lb.this.zone_id }
output "alb_arn" { value = aws_lb.this.arn }
output "order_bot_target_group_arn" { value = aws_lb_target_group.svc.arn }
output "order_bot_mgmt_target_group_arn" { value = aws_lb_target_group.mgmt.arn }
