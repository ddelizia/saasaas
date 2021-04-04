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

func TestGetShared(test *testing.T) {
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
		result := &ExampleDataShared{}
		err := dyDb.Repo.Get(context.Background(), t.NewString(id), result)

		// Then
		assert.Nil(test, err)
		assert.Equal(test, result.DataInt, 1)
		assert.Equal(test, result.DataString, id)

	})

}
