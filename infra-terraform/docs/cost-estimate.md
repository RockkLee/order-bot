# Cost estimate (ALB + ECS active 8 hours/day)

This is a **rough estimate** based on the current Terraform defaults:

- 2 ECS Fargate services (`order-bot-svc`, `order-bot-mgmt-svc`)
- each service runs `desired_count = 1` when active
- task size defaults to `0.25 vCPU` and `0.5 GB memory` per task
- scheduled scaling is 10:00 to 18:00 (8 hours/day)
- estimate assumes 30 days/month and low traffic (about 1 ALB LCU while active)

## Assumptions used

- Active hours per month: `8 × 30 = 240` hours
- ECS task-hours per month: `2 tasks × 240 = 480 task-hours`

## Formula

`monthly_cost ≈ ECS(Fargate vCPU + memory) + ALB(hourly + LCU) + small fixed services`

Where:

- `ECS = 480 × (0.25 × vCPU_price_per_hr + 0.5 × GB_price_per_hr)`
- `ALB = 240 × (ALB_hourly_price + LCU_hourly_price)`

## Ballpark numbers

Using typical Linux/x86 on-demand Fargate and ALB rates in common regions, this lands around:

- **ECS (2 tasks, 8h/day):** about **$6–$9 / month**
- **ALB (8h/day + ~1 LCU):** about **$8–$12 / month**
- **CloudWatch logs, Route53 records, ECR storage, CloudFront/S3 (light usage):** about **$2–$10 / month**

### Estimated total

- **~$16 to $31 / month** for light usage.

## Important caveat about ALB runtime

ALB does not have a native "stop/start" like EC2. In practice, ALB hourly billing usually continues while the ALB exists.

So there are two practical scenarios:

1. **If ALB truly exists only 8h/day** (created/destroyed by automation): estimate above.
2. **If ALB exists 24/7** (most common): add roughly **3x ALB hourly component**, and monthly total is more likely around **$30–$55+** depending on traffic and region.

For production budgeting, plug exact values into AWS Pricing Calculator with your region and expected traffic/profile.
