# Lambda mock server

This is mock json requests service backed by AWS Lambda, Api Gateway, Dynamodb

## To Deploy

Make sure your AWS credentials are exported. Example:

```
export AWS_ACCESS_KEY_ID=the-key && export AWS_SECRET_ACCESS_KEY=the-secret

make tf-deploy
```

## Local Development

Requires golang >= 1.23

```
make dev
```

Available at: http://localhost:8000

### Using Docker compose

```
make up

docker-compose exec api /bin/bash

make dev
```
