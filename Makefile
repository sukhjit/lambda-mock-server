BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64

LDFLAGS=-ldflags "-s -w"

dev:
	$(BUILD_ENV) go build -o bin/main main.go

clean:
	rm -rf .serverless ./bin

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/main main.go

deploy: clean build
	./node_modules/.bin/sls deploy --verbose
