output "alb_dns_name" {
  # In Terraform, a counted module becomes a list of instances. So you must index it.
  value = var.enable_alb ? module.alb[0].alb_dns_name : null
}

output "ecs_cluster_name" {
  value = module.ecs.cluster_name
}
