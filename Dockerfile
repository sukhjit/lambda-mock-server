FROM golang:1.14-alpine

RUN apk --no-cache add --update wget curl git bash ca-certificates make \
	&& go get github.com/githubnemo/CompileDaemon

WORKDIR /app

EXPOSE 8000

ENTRYPOINT CompileDaemon -build="make dev" -command="./bin/main"
