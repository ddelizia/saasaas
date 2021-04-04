package dynamodbsaas_test

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/db/dynamodbsaas"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type AwsDynamodbSetup struct {
	AwsConfig *aws.Config
	Session   *session.Session
	TableName string
	Db        *dynamodb.DynamoDB
	Repo      *dynamodbsaas.Db
}

type ExampleDataShared struct {
	dynamodbsaas.DynamoSaasModel
	ID         string
	DataInt    int
	DataString string
}

type ExampleDataAccount struct {
	dynamodbsaas.DynamoSaasModel
	db.AccountAware
	ID         string
	DataInt    int
	DataString string
}

type ExampleDataUser struct {
	dynamodbsaas.DynamoSaasModel
	db.UserAware
	ID         string
	DataInt    int
	DataString string
}

func SetupDynamodb(t *testing.T) *AwsDynamodbSetup {
	awsConfig := &aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials("DUMMY", "DUMMY", "token"),
		Endpoint:    aws.String("http://localhost:8000"),
	}
	session := session.Must(session.NewSession(awsConfig))
	tableName := "DyTestTable-" + uuid.New().String()
	db := dynamodb.New(session, awsConfig)

	_, err := db.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("pk"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	})

	assert.Nil(t, err, "Error creating the dynamodb table")

	repoSharedDf := map[string]dynamodbsaas.DynamicAttrFunc{}
	dynamodbsaas.StdDFWrapperSharedModel("#SHARED", "#DATA", repoSharedDf)

	repoAccountDf := map[string]dynamodbsaas.DynamicAttrFunc{}
	dynamodbsaas.StdDFWrapperAccountModel("#DATA_REPO_ACCOUNT", repoAccountDf)

	repoUserDf := map[string]dynamodbsaas.DynamicAttrFunc{}
	dynamodbsaas.StdDFWrapperAccountModel("#DATA_REPO_USER", repoUserDf)

	pk := &dynamodbsaas.Index{
		PkField: "pk",
		SkField: "sk",
	}

	repo := &dynamodbsaas.Db{
		DynamoDb:  db,
		TableName: tableName,
		EntityDescriptorMapping: map[reflect.Type]*dynamodbsaas.EntityDescriptor{
			reflect.TypeOf(&ExampleDataShared{}): {
				DynamicFields: repoSharedDf,
				IDBuilder:     dynamodbsaas.StdIDBuilderSharedModel("#SHARED", "#DATA"),
				PK:            pk,
			},
			reflect.TypeOf(&ExampleDataAccount{}): {
				DynamicFields: repoAccountDf,
				IDBuilder:     dynamodbsaas.StdIDBuilderAccountModel("#DATA_REPO_ACCOUNT"),
				PK:            pk,
			},
			reflect.TypeOf(&ExampleDataUser{}): {
				DynamicFields: repoUserDf,
				IDBuilder:     dynamodbsaas.StdIDBuilderUserModel("#DATA_REPO_USER"),
				PK:            pk,
			},
		},
	}

	return &AwsDynamodbSetup{
		AwsConfig: awsConfig,
		Session:   session,
		TableName: tableName,
		Db:        db,
		Repo:      repo,
	}
}
