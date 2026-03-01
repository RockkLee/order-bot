# Wiring AWS credentials for `terraform init`

This repo uses an S3 backend for Terraform state. The S3 backend is configured in:
- `environments/global/providers.tf`
- `environments/prod/providers.tf`

Important: the backend initializes *before* Terraform variables are loaded, so it **cannot** read
`aws_profile` from `terraform.tfvars`. You must provide backend credentials via environment variables
or a backend config file at init time.

## Prerequisites

1) AWS CLI v2 installed.
2) An AWS SSO profile configured and logged in.
   - aws configure sso
      - ```bash
        SSO session name (Recommended): # type whatever you want
        SSO start URL [None]: # get it from the AWS Web page: "IAM Identity Center > Settings > Identity source > AWS access portal URLs > IPv4-only > [the actual url!!!]>"
        SSO region [None]: ap-northeast-1
        ```

## Option A: Use environment variables (recommended)

From the environment directory (e.g. `environments/global`):

```bash
AWS_PROFILE=ordebot-profile AWS_SDK_LOAD_CONFIG=1 terraform init
```

Repeat for `environments/prod`.

## Option B: Use a backend config file

Create a local `backend.hcl` (do **not** commit it). You can copy from:
- `environments/global/backend.hcl.example`
- `environments/prod/backend.hcl.example`

Then run:

```bash
terraform init -backend-config=backend.hcl
```

