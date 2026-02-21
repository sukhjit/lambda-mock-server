resource "aws_apigatewayv2_api" "api" {
  name          = local.app_name
  protocol_type = "HTTP"
  tags          = local.tags
}

resource "aws_cloudwatch_log_group" "gateway_stage" {
  name              = "API-Gateway-Execution-Logs_${aws_apigatewayv2_api.api.id}/${var.environment}"
  retention_in_days = 90
  tags              = local.tags
}

resource "aws_apigatewayv2_integration" "api" {
  api_id               = aws_apigatewayv2_api.api.id
  integration_type     = "AWS_PROXY"
  integration_method   = "POST"
  integration_uri      = aws_lambda_function.rest_api.arn
  passthrough_behavior = "WHEN_NO_MATCH"
}

resource "aws_apigatewayv2_route" "api" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "ANY /{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.api.id}"
}

resource "aws_lambda_permission" "lambda_api" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  principal     = "apigateway.amazonaws.com"
  function_name = aws_lambda_function.rest_api.function_name
  source_arn    = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}

resource "aws_apigatewayv2_stage" "site" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = "$default"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.gateway_stage.arn
    format          = jsonencode(jsondecode(trimspace(file("${path.module}/template/api-gw-access-log-format.json"))))
  }

  route_settings {
    route_key                = aws_apigatewayv2_route.api.route_key
    throttling_rate_limit    = 1000
    throttling_burst_limit   = 250
    detailed_metrics_enabled = true
  }

  tags = merge(local.tags, {
    Name = local.app_name
  })
}

output "api_url" {
  value = aws_apigatewayv2_api.api.api_endpoint
}
