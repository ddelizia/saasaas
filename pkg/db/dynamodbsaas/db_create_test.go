package dynamodbsaas_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ddelizia/saasaas/pkg/t"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateShared(test *testing.T) {
	dyDb := SetupDynamodb(test)

	test.Run("it should return the data correctly", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		example := &ExampleDataShared{
			ID:         id,
			DataString: id,
			DataInt:    1,
		}

		// When
		err := dyDb.Repo.Create(context.Background(), t.NewString(example.DataString), example)

		// Then
		assert.Nil(test, err)

		get, err := dyDb.Db.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(dyDb.TableName),
			Key: map[string]*dynamodb.AttributeValue{
				"pk": {
					S: aws.String("#SHARED"),
				},
				"sk": {
					S: aws.String("#DATA#" + id),
				},
			},
		})
		result := &ExampleDataShared{}
		dynamodbattribute.UnmarshalMap(get.Item, result)
		assert.Equal(test, result.DataInt, 1)
		assert.Equal(test, result.DataString, id)

	})

	test.Run("it should fail when I try to insert 2 times the same pk", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		example := &ExampleDataShared{
			ID:         id,
			DataString: id,
			DataInt:    1,
		}

		// When
		dyDb.Repo.Create(context.Background(), t.NewString(example.DataString), example)
		err := dyDb.Repo.Create(context.Background(), t.NewString(example.DataString), example)

		// Then
		assert.NotNil(test, err)
	})
}
