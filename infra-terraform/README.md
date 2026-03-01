# infra-terraform

Terraform code for deploying the order-bot system shown in the architecture diagram.

## Directory layout

- `docs/`: architecture and operational notes.
- `docs/cost-estimate.md`: rough monthly cost estimate and pricing caveats.
- `environments/global`: global resources (ECR, frontend S3 + CloudFront + Route53).
- `environments/prod`: production runtime resources (ALB, ECS, SG, scheduled scaling).
- `modules/`: reusable Terraform modules.

## Modules

- `alb`: public ALB with host-based routing for `order-bot-svc` and `order-bot-mgmt-svc`.
- `security_group`: two SGs (`ALB`, shared app SG).
- `ecr`: image repositories and lifecycle policies.
- `ecs`: ECS cluster, task definitions, services, logs, and autoscaling targets.
- `s3`: frontend hosting bucket + CloudFront + Route53 alias.
- `schedule`: scale up/down schedules (10:00-18:00).

## Apply order

1. `environments/global`
2. `environments/prod`

## Notes

- Backend services are deployed on ECS Fargate behind ALB.
- A shared app SG is used for ECS services and PostgreSQL EC2 access to match the requested "one security" model.
- Autoscaling targets support scale-to-zero out of business hours.

## Commands
- You can disable the ALB (and its Route53 aliases) by applying `environments/prod` with `enable_alb = false`
```bash
terraform apply -var=enable_alb=false
```
