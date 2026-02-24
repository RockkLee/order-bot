# Terraform Start Flow and Resource Inventory

This document summarizes what is created from all `main.tf` and `outputs.tf` files under `infra-terraform`, and the high-level startup/runtime flow between resources.

## 1) Root environment composition

There are two environment entry points:

- `environments/global`: shared/global resources (ECR + frontend hosting).
- `environments/prod`: runtime resources for production API services (Security Groups + ALB + ECS + scheduled scaling + DNS).

## 2) Mermaid: end-to-end start flow

```mermaid
flowchart TD
  A[Terraform apply]

  subgraph G[Global environment]
    G1[module.ecr]
    G2[module.frontend (S3 + CloudFront + Route53)]
    G3[(Outputs: ecr_repository_urls, frontend_bucket_name, frontend_cloudfront_domain)]
    G1 --> G3
    G2 --> G3
  end

  subgraph P[Prod environment]
    P1["module.security_group<br/>creates ALB SG + APP SG"]
    P2["module.alb<br/>creates ALB, listeners, host rules, target groups"]
    P3["Route53 aliases<br/>orderbot + orderbot_mgmt -> ALB"]
    P4["module.ecs<br/>cluster, roles, task defs, services, autoscaling targets"]
    P5["module.schedule<br/>cron scale up/down for both ECS services"]
    P6[(Outputs: alb_dns_name, ecs_cluster_name)]

    P1 --> P2
    P2 --> P3
    P1 --> P4
    P2 --> P4
    P4 --> P5
    P2 --> P6
    P4 --> P6
  end

  B[Users / Clients] -->|HTTPS| C[Route53 records]
  C --> D[ALB 443 listener]
  D -->|Host: orderbot domain| E[order-bot target group]
  D -->|Host: orderbot-mgmt domain| F[order-bot-mgmt target group]
  E --> G4[ECS service: order-bot-svc]
  F --> G5[ECS service: order-bot-mgmt-svc]

  H[Frontend users] --> I[Frontend Route53 alias]
  I --> J[CloudFront]
  J --> K[S3 frontend bucket via OAC]

  L[CI/CD image push] --> M[ECR repositories]
  M --> N[ECS task definitions pull images]
  N --> G4
  N --> G5
```

## 3) Resource inventory by module

### `environments/global`

#### `module.ecr`
Creates one ECR repository per configured name, plus a lifecycle policy per repository.

- `aws_ecr_repository.this` (for_each)
- `aws_ecr_lifecycle_policy.this` (for_each)

Outputs exposed from module and environment:

- `repository_urls` (module output)
- `ecr_repository_urls` (environment output)

#### `module.frontend` (`modules/s3`)
Creates static hosting and CDN chain.

- `aws_s3_bucket.frontend`
- `aws_s3_bucket_public_access_block.frontend`
- `aws_cloudfront_origin_access_control.this`
- `aws_cloudfront_distribution.this`
- `aws_s3_bucket_policy.frontend` (allow CloudFront read)
- `aws_route53_record.frontend_alias`

Outputs exposed from module and environment:

- `bucket_name` -> `frontend_bucket_name`
- `cloudfront_domain_name` -> `frontend_cloudfront_domain`

### `environments/prod`

#### `module.security_group`
- `aws_security_group.alb` (ingress 80/443 from allowed CIDRs)
- `aws_security_group.app` (ingress app port from ALB SG + self 5432)

Outputs:

- `alb_security_group_id`
- `app_security_group_id`

#### `module.alb`
- `aws_lb.this`
- `aws_lb_target_group.orderbot`
- `aws_lb_target_group.orderbot_mgmt`
- `aws_lb_listener.https` (TLS termination + default 404)
- `aws_lb_listener.http_redirect` (80 -> 443)
- `aws_lb_listener_rule.orderbot_mgmt_by_host`
- `aws_lb_listener_rule.orderbot_by_host`

Outputs used by prod root:

- `alb_dns_name`, `alb_zone_id`, `alb_arn`
- `order_bot_target_group_arn`, `order_bot_mgmt_target_group_arn`

#### Prod Route53 records
At root prod level:

- `aws_route53_record.orderbot_mgmt_alias`
- `aws_route53_record.orderbot_alias`

Both alias to ALB DNS/zone outputs from `module.alb`.

#### `module.ecs`
Creates ECS runtime and service-level autoscaling targets.

- Data sources:
  - `aws_region.current`
  - `aws_caller_identity.current`
- Logs:
  - `aws_cloudwatch_log_group.svc`
  - `aws_cloudwatch_log_group.mgmt`
- Compute/control plane:
  - `aws_ecs_cluster.this`
  - `aws_iam_role.task_execution`
  - `aws_iam_role_policy_attachment.task_exec_default`
  - `aws_iam_role.task_role`
  - `aws_ecs_task_definition.svc`
  - `aws_ecs_task_definition.mgmt`
  - `aws_ecs_service.svc`
  - `aws_ecs_service.mgmt`
  - `aws_appautoscaling_target.svc`
  - `aws_appautoscaling_target.mgmt`

Outputs:

- `cluster_name`
- `order_bot_autoscaling_resource_id`
- `order_bot_mgmt_autoscaling_resource_id`

#### `module.schedule`
Creates scheduled actions bound to ECS service autoscaling targets.

- `aws_appautoscaling_scheduled_action.scale_up` (for_each services)
- `aws_appautoscaling_scheduled_action.scale_down` (for_each services)

#### Prod environment outputs
- `alb_dns_name`
- `ecs_cluster_name`

## 4) Practical startup sequence

1. Apply `global` to provision image registry and frontend delivery stack.
2. Build/push service images to the ECR repos output from `global`.
3. Apply `prod` to provision networking, ALB routing, ECS cluster/services, and DNS aliases.
4. ECS services start tasks from configured image URIs and register into ALB target groups.
5. Scheduled autoscaling actions toggle desired counts on cron windows.
