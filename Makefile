GO=go
GO_BUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags "-s -w"
pwdDir=$(shell pwd)
distDir=$(pwdDir)/dist
lintVersion=v1.45.2
lintURL=https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh
linter=$(distDir)/golangci-lint

# dev
lint-install:
	@mkdir -p $(distDir)
	@[ -f "$(linter)" ] || curl -sSfL $(lintURL) | sh -s -- -b $(distDir) $(lintVersion)

lint: lint-install
	@$(linter) run -v --skip-dirs='node_modules' ./...

compile-daemon-install:
	@[ -f "$(distDir)/CompileDaemon" ] || GOBIN=$(distDir) go install github.com/githubnemo/CompileDaemon@v1.4.0

dev: compile-daemon-install
	@$(distDir)/CompileDaemon -build="go build -o dist/mock-server main.go" -command="./dist/mock-server"

#  build
clean:
	rm -rf .serverless ./bin

build:
	$(GO_BUILD) -o bin/main main.go

deploy: clean build
	./node_modules/.bin/sls deploy --verbose
