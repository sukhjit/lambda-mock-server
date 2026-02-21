resource "aws_cloudwatch_log_group" "lambda_cloudwatch" {
  name              = "/aws/lambda/${local.lambda_api_name}"
  retention_in_days = local.log_rentention
  tags              = local.tags
}

resource "aws_iam_role" "app_role" {
  name = "${local.app_name}-app-role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF

  tags = merge(local.tags, {
    Name = "${local.app_name}-api"
  })
}

resource "aws_iam_policy" "app_policy" {
  name   = "${local.app_name}-app-policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": [
        "${aws_cloudwatch_log_group.lambda_cloudwatch.arn}:*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "xray:PutTelemetryRecords",
        "xray:PutTraceSegments"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "app_policy" {
  role       = aws_iam_role.app_role.name
  policy_arn = aws_iam_policy.app_policy.arn
}

data "archive_file" "api" {
  type        = "zip"
  source_file = "${path.root}/../../dist/bin/api/bootstrap"
  output_path = "${path.root}/../../dist/zip/api.zip"
}

resource "aws_lambda_function" "rest_api" {
  function_name    = local.lambda_api_name
  filename         = data.archive_file.api.output_path
  description      = "Lmabda Mock Server - ${var.environment}"
  handler          = "bootstrap"
  source_code_hash = data.archive_file.api.output_base64sha256
  role             = aws_iam_role.app_role.arn
  runtime          = "provided.al2023"
  memory_size      = 256
  timeout          = 30
  architectures    = ["arm64"]

  tags = merge(
    local.tags,
    {
      purpose = "Lambda function for API."
    }
  )

  environment {
    variables = {
      ENVIRONMENT = var.environment
      GIN_MODE    = "release"
    }
  }
}
