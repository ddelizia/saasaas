package dynamodbsaas_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/ddelizia/saasaas/pkg/t"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteShared(test *testing.T) {
	dyDb := SetupDynamodb(test)

	test.Run("it should return correctly the data", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		dyDb.Db.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(dyDb.TableName),
			Item: map[string]*dynamodb.AttributeValue{
				"pk": {
					S: aws.String("#SHARED"),
				},
				"sk": {
					S: aws.String("#DATA#" + id),
				},
				"ID": {
					S: aws.String(id),
				},
				"DataInt": {
					N: aws.String("1"),
				},
				"DataString": {
					S: aws.String(id),
				},
			},
		})

		// When
		err := dyDb.Repo.Delete(context.Background(), t.NewString(id), &ExampleDataShared{})

		// Then
		assert.Nil(test, err)

		get, err := dyDb.Db.GetItem(&dynamodb.GetItemInput{
			TableName: &dyDb.TableName,
			Key: map[string]*dynamodb.AttributeValue{
				"pk": {
					S: aws.String("#SHARED"),
				},
				"sk": {
					S: aws.String("#DATA#" + id),
				},
			},
		})
		assert.Nil(test, err)
		assert.Nil(test, get.Item)
	})

}
