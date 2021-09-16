GO=go
GO_BUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w"

clean:
	rm -rf .serverless ./bin

build:
	$(GO_BUILD) -o bin/main main.go

deploy: clean build
	./node_modules/.bin/sls deploy --verbose
