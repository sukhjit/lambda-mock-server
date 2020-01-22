clean:
	rm -rf .serverless ./bin

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/main main.go

deploy: clean build
	./node_modules/.bin/sls deploy --verbose
