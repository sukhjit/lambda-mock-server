package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/sukhjit/lambda-mock-server/handler"
)

var (
	ginLambda *ginadapter.GinLambda
	router    *gin.Engine
)

func main() {
	isLambda := len(os.Getenv("AWS_LAMBDA_FUNCTION_NAME")) > 0

	router = handler.New()

	if isLambda {
		lambda.Start(lambdaHandler)
	} else {
		if err := router.Run(":8000"); err != nil {
			log.Fatal(err)
		}

	}
}

func lambdaHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		ginLambda = ginadapter.New(router)
	}

	return ginLambda.Proxy(req)
}
