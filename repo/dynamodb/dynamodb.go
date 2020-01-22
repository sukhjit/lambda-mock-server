package dynamodb

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sukhjit/lambda-mock-server/model"
	"github.com/sukhjit/lambda-mock-server/repo"
)

const (
	budgetTableName = "rds-01-BudgetDynamoTable-FCG2RILOQ5EZ"
	awsRegion       = "ap-southeast-2"
)

type document struct {
	db            *dynamodb.DynamoDB
	documentTable string
}

// New func
func New(awsRegion, docTable string) repo.Document {
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	return &document{
		db:            dynamodb.New(session),
		documentTable: docTable,
	}
}

func (d *document) Get(id string) (*model.Document, error) {
	log.Println("Table", d.documentTable)
	result, err := d.db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.documentTable),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	item := model.Document{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (d *document) Add(document *model.Document) error {
	_, err := d.db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.documentTable),
		Item: map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(document.ID),
			},
			"body": &dynamodb.AttributeValue{
				S: aws.String(fmt.Sprintf("%v", document.Body)),
			},
			"date": &dynamodb.AttributeValue{
				S: aws.String(time.Now().Format("2006-01-02 15:04:05")),
			},
		},
	})

	return err
}

func (d *document) Close() error {
	return nil
}
