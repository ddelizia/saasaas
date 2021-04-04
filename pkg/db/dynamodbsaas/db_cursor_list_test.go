package dynamodbsaas_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCursorListShared(t *testing.T) {
	dyDb := SetupDynamodb(t)

	ids := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	for i, id := range ids {
		dyDb.Db.PutItem(&dynamodb.PutItemInput{
			TableName: &dyDb.TableName,
			Item: map[string]*dynamodb.AttributeValue{
				"pk": {
					S: aws.String("#SHARED"),
				},
				"sk": {
					S: aws.String("#DATA#" + strconv.Itoa(i)),
				},
				"ID": {
					S: aws.String(strconv.Itoa(i)),
				},
				"DataInt": {
					N: aws.String(strconv.Itoa(i)),
				},
				"DataString": {
					S: aws.String(id),
				},
			},
		})
	}

	t.Run("it should return the first 2 elements only", func(t *testing.T) {
		// Given
		var elements []*ExampleDataShared

		// When
		l, err := dyDb.Repo.CursorList(context.Background(), "", 2, &elements)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, len(elements), 2)

		for _, element := range elements {
			assert.Contains(t, ids, element.DataString)
		}
		assert.NotNil(t, l)
		assert.NotNil(t, l.LastKey)
	})

	t.Run("it should get the first 2 pages", func(t *testing.T) {
		// Given
		var elements1 []*ExampleDataShared
		l1, err := dyDb.Repo.CursorList(context.Background(), "", 2, &elements1)
		assert.Nil(t, err, "the first page must exists")

		// When
		var elements2 []*ExampleDataShared
		l2, err := dyDb.Repo.CursorList(context.Background(), l1.LastKey, 2, &elements2)

		// Then
		assert.Nil(t, err)
		assert.Equal(t, len(elements2), 2)

		for _, element := range elements2 {
			assert.Contains(t, ids, element.DataString)
		}
		assert.NotNil(t, l2)
		assert.NotNil(t, l2.LastKey)
	})

	t.Run("it should have empty next page when I reach the end of the list", func(t *testing.T) {
		// Given
		var elements1 []*ExampleDataShared
		l1, err := dyDb.Repo.CursorList(context.Background(), "", 5, &elements1)

		// When
		var elements2 []*ExampleDataShared
		l2, err := dyDb.Repo.CursorList(context.Background(), l1.LastKey, 5, &elements2)

		// Then
		assert.Nil(t, err)
		assert.Empty(t, elements2)
		assert.Equal(t, "", l2.LastKey)
	})

}
