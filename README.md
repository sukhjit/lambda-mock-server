# Lambda mock server

This is mock json requests service backed by AWS Lambda, Api Gateway, Dynamodb

## To Deploy

Make sure your AWS credentials are exported. Example:
```
export AWS_ACCESS_KEY_ID=the-key && export AWS_SECRET_ACCESS_KEY=the-secret

npm install
make deploy
```

## Local Development

Requires npm, golang >= 1.13 and dynamodb table created in AWS
```
export AWS_ACCESS_KEY_ID=the-key && export AWS_SECRET_ACCESS_KEY=the-secret
export AWS_REGION="us-east-1"
export DOCUMENT_TABLE_NAME="dynamodb-table-name-in-aws"

go run main.go
```
