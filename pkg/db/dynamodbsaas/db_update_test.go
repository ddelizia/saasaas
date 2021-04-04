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

func TestUpdateShared(test *testing.T) {
	dyDb := SetupDynamodb(test)

	test.Run("it should fail when I try to update something that does not exists", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		example := &ExampleDataShared{
			ID:         id,
			DataString: id,
			DataInt:    1,
		}

		// When
		err := dyDb.Repo.Update(context.Background(), t.NewString(example.ID), example)

		// Then
		assert.NotNil(test, err)
	})

	test.Run("it should update the data correctly", func(test *testing.T) {
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
		example := &ExampleDataShared{
			ID:         id,
			DataString: "Hello",
			DataInt:    2,
		}
		err := dyDb.Repo.Update(context.Background(), t.NewString(id), example)

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

		assert.Equal(test, result.DataString, "Hello")
		assert.Equal(test, result.DataInt, 2)

	})
}
