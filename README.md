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

Requires npm, golang >= 1.13

Available at: http://localhost:8000

```
cp .env.dist .env

# hot reload
cicd/bin/local-dev.sh
```

### Using Docker compose

```
docker-compose up

docker-compose exec api bash

cp .env.dist .env

# hot reload
cicd/bin/local-dev.sh
```

## Serving

```
https://rbait02cw7.execute-api.ap-southeast-2.amazonaws.com/prod
```
