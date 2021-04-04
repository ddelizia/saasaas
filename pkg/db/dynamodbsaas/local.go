package dynamodbsaas

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func DynamoDbInstance() *dynamodb.DynamoDB {
	env := os.Getenv("ENV")
	awsConfig := &aws.Config{}
	if env == "" || env == "local" {
		awsConfig = &aws.Config{
			Region:      aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials("DUMMY", "DUMMY", "token"),
			Endpoint:    aws.String("http://localhost:8000"),
		}
	}
	session := session.Must(session.NewSession(awsConfig))
	return dynamodb.New(session, awsConfig)
}

func TableNameWithEnv(tableName string) string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}
	return fmt.Sprint(tableName, "-", env)
}
