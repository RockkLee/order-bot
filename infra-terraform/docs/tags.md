# Tags

Tags are key/value metadata attached to AWS resources. They are used for organization, cost allocation, searchability, and automation (e.g., policies, scripts, or resource groups).

## Tags in this repo

The common tag set is defined as a Terraform local called `locals.tags` and passed into modules so that resources inherit the same metadata.

### Definitions

- `environments/global/main.tf`
  - `Project = "order-bot"`
  - `Environment = var.environment`
  - `ManagedBy = "terraform"`

- `environments/prod/main.tf`
  - `Project = "order-bot"`
  - `Environment = "prod"`
  - `ManagedBy = "terraform"`

### Where they are used

- `environments/global/main.tf`: modules `ecr`, `frontend`
- `environments/prod/main.tf`: modules `security_group`, `alb`, `ecs`

Note: `aws_route53_record` resources do not support tags, so they are untagged.

## Using tags in the AWS Console

The UI varies by service, but the following patterns are consistent.

### Filter/search within a service

Most service list pages (EC2, ECS, ECR, S3, ALB, etc.) provide a filter bar. Use the filter option for **Tags** or **Tag key/value** to narrow the list to resources with a specific tag (e.g., `Project=order-bot` or `Environment=prod`).

### Tag Editor (cross-service search)

Use **Tag Editor** in the **Resource Groups & Tag Editor** service to:

- Search across multiple services by tag key/value
- Bulk-add or edit tags
- Export the results

### Resource Groups

Create a **Resource Group** based on tag rules (e.g., `Project=order-bot AND Environment=prod`) to manage and view related resources together.

## Recommended tag usage

- Use `Project` to identify ownership or product
- Use `Environment` to separate prod/dev/test
- Use `ManagedBy` to distinguish Terraform-managed resources
