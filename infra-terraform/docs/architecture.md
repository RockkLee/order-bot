# Architecture mapping

This Terraform implementation follows the provided Mermaid diagram and the discovered app behavior:

- **Frontend (`order-bot-frontend`)** expects:
  - `VITE_ORDER_BOT_BASE_PATH` for chat API
  - `VITE_ORDER_BOT_MGMT_BASE_PATH` for management API
- **Chat service (`order-bot-svc`)**:
  - FastAPI app with `/chat` endpoints and menu lookups.
  - Runtime configured via env vars (`HOST`, `PORT`, `DATABASE_URL`, etc.).
- **Management service (`order-bot-mgmt-svc`)**:
  - Gin app under `/orderbotmgmt` plus `/health`.
  - Runtime configured via env vars (`ADDRESS`, `PORT`, `BLUEPRINT_DB_*`, etc.).

## Diagram-to-resource translation

- Route53 -> ALB for API domains (`chat`, `mgmt`) in `environments/prod`.
- Route53 -> CloudFront -> S3 for frontend domain in `environments/global`.
- ALB across two public subnets/AZs (input-driven).
- ECS services and PostgreSQL EC2 are assumed to share one app SG.
- ECS task definitions use ECR images and ship logs to CloudWatch Logs.
- Scheduled autoscaling sets desired capacity to 1 during 10:00-18:00 and allows 0 outside that window.

```mermaid
graph TD
    subgraph "Networking"
        R53[AWS Route 53] --> ALB[AWS Application Loadbalancer]
        ALB --> AZa
        ALB --> AZb
        AZa --> SG[Security Groups:<br>one for ALB, on for others]
        AZb --> SG
    end

    subgraph "Application Layer"
        SG --putting everthing in one security--> ECS_Svc[ECS Services]
        ECS_Svc --> ECS[AWS ECS Cluster]
        ECS_Svc -- Uses --> ECS_TD[ECS Task Definitions]
        ECS_TD -- Uses Images from --> ECR
        ECS_Svc -- Logs to --> CWL[AWS CloudWatch Logs]
    end

    subgraph "Data Layer"
        SG --putting everthing in one security--> EC2[EC2 with a PostgreSQL DB running on Docker]
    end

    subgraph "Frontend Hosting"
        R53 --> CF[AWS CloudFront]
    end

    subgraph "Scheduled scaling 10:00-18:00"
        ECS -- 1.Auto Scaling to zero<br>2.supports no downtime rolling but by default is one instance if no need to roll a version --> AAS[AWS Application Auto Scaling]
        ALB --> Unknown
        CF --> Unknown
    end
```

```mermaid
flowchart TD
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

## Gaps and assumptions

- PostgreSQL EC2 instance provisioning is intentionally external to keep module list aligned with the request; DB host is injected via environment variables.
- Existing VPC/subnets/hosted zone/certificates are expected as inputs.
