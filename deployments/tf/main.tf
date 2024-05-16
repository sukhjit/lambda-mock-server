terraform {
  backend "s3" {
    bucket         = "terraform-state-files-01"
    region         = "ap-southeast-2"
    encrypt        = true
    dynamodb_table = "terraform-lock-01"
  }

  required_providers {
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2"
    }

    aws = {
      source  = "hashicorp/aws"
      version = "~> 5"
    }
  }
}

provider "aws" {
  region = "ap-southeast-2"
}
