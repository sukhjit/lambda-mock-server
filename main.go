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
	ginLambda     *ginadapter.GinLambda
	isLambda      bool
	router        *gin.Engine
	awsRegion     string
	documentTable string
)

func initEnv() {
	_ = godotenv.Load()

	awsRegion = getEnv("AWS_REGION", "ap-southeast-2")

	documentTable = getEnv("DOCUMENT_TABLE_NAME", "lambda-mock-server-document")

	isLambda = false
	if os.Getenv("WEB") == "" {
		isLambda = true
	}
}

func main() {
	initEnv()

	router = handler.New(awsRegion, documentTable)

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

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
