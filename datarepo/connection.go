package datarepo

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var configuration *aws.Config
var dynamodbConnection *dynamodb.DynamoDB
var tableName *string
var logtableName *string

func init() {
	DynamoDBService() // Preload this
}

//DynamoDBService returns a properly configured dynamodb service
func DynamoDBService() *dynamodb.DynamoDB {
	if configuration == nil || dynamodbConnection == nil {
		var err error
		configuration, err = loadConfiguration()
		if err != nil {
			panic("Failed to load AWS DynamoDB configuration: " + err.Error())
		}
		dynamodbConnection = dynamodb.New(*configuration)
	}

	return dynamodbConnection
}

func DataTableName() string {
	if tableName == nil {
		value := os.Getenv("TABLE_NAME")
		tableName = &value
	}
	return *tableName
}

func loadConfiguration() (*aws.Config, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
