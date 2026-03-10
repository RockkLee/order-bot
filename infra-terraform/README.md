# infra-terraform

## Commands
- You can disable the ALB (and its Route53 aliases) by applying `environments/prod` with `enable_alb = false`
```bash
terraform apply -var=enable_alb=false
```
