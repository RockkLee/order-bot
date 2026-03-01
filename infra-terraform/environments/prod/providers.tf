terraform {
  required_version = ">= 1.6.0"

  backend "s3" {
    bucket = "order-bot-terraform-state"
    key    = "prod/terraform.tfstate"
    region = "ap-northeast-2"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}
