package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sukhjit/lambda-mock-server/handler"
)

var (
	ginLambda *ginadapter.GinLambda
	isLambda  bool
	router    *gin.Engine
)

func initEnv() {
	_ = godotenv.Load()

	isLambda = false
	if os.Getenv("WEB") == "" {
		isLambda = true
	}
}

func main() {
	initEnv()

	router = handler.New()

	if isLambda {
		lambda.Start(lambdaHandler)
	} else {
		router.Run(":8000")
	}
}

func lambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		ginLambda = ginadapter.New(router)
	}

	return ginLambda.Proxy(req)
}
