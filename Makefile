applicationName=lambda-mock-server
goBuild=CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -trimpath -ldflags='-w -s -extldflags "-static"'
pwdDir=$(shell pwd)
distDir=$(pwdDir)/dist
binaryOutputDir=$(distDir)/bin
lintVersion=v1.58.1
lintURL=https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh
linter=$(distDir)/golangci-lint
terraformer=$(distDir)/terraform
terraformVersion=1.1.9
terraformURL="https://releases.hashicorp.com/terraform/$(terraformVersion)/terraform_$(terraformVersion)_linux_amd64.zip"

# docker
up:
	@docker-compose up --build -d api

down:
	@docker-compose down -v

# dev
lint-install:
	@mkdir -p $(distDir)
	@[ -f "$(linter)" ] || curl -sSfL $(lintURL) | sh -s -- -b $(distDir) $(lintVersion)

lint: lint-install
	@$(linter) run -c .golangci.yml -v ./...

compile-daemon-install:
	@[ -f "$(distDir)/CompileDaemon" ] || GOBIN=$(distDir) go install github.com/githubnemo/CompileDaemon@v1.4.0

dev: compile-daemon-install
	@$(distDir)/CompileDaemon -build="go build -o dist/mock-server main.go" -command="./dist/mock-server"

build:
	@mkdir -p $(binaryOutputDir)/api
	$(goBuild) -o $(binaryOutputDir)/api/bootstrap main.go

# tf
install-pre-req:
	apt-get update && apt-get install -y curl zip unzip

tf-install: install-pre-req
	@[ -f "$(terraformer)" ] \
		|| (curl -sSL $(terraformURL) -o /tmp/terraform.zip \
			&& unzip /tmp/terraform.zip -d $(distDir)/ \
			&& rm /tmp/terraform.zip && $(terraformer) version)

tf-version: tf-install
	@$(terraformer) version

tf-clean:
	@rm -rf deployments/tf/.terraform deployments/tf/tfplan deployments/tf/.terraform.lock.hcl

tf-deploy: tf-version tf-clean build
	cd deployments/tf \
	&& $(terraformer) init \
		-backend-config="key=$(applicationName)/prod/terraform/state.tf" \
		-input=false \
	&& $(terraformer) validate \
	&& $(terraformer) plan \
		-out=tfplan \
	&& $(terraformer) apply tfplan
	make tf-clean
