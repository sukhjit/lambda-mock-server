variable "application" {
  type    = string
  default = "lambda-mock-server"
}

variable "environment" {
  type    = string
  default = "prod"
}

locals {
  log_rentention  = 90
  app_name        = "${var.application}-${var.environment}"
  lambda_api_name = "${local.app_name}-lambda-api"
  api_gw_name     = local.app_name

  tags = {
    application = var.application
    provision   = "Terraform"
    environment = var.environment
  }
}
