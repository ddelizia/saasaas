package mongosaas_test

import (
	"context"
	"testing"

	"github.com/ddelizia/saasaas/pkg/t"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateShared(test *testing.T) {
	mgmTestData := SetupMongo(test)

	test.Run("it should return the data correctly", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		example := &ExampleDataShared{
			DataString: t.NewString(id),
			DataInt:    t.NewInt64(1),
		}

		// When
		err := mgmTestData.Db.Create(context.Background(), t.NewString(id), example)

		// Then
		assert.Nil(test, err)

		result := ExampleDataShared{}

		assert.Equal(test, result.DataInt, t.NewInt64(1))
		assert.Equal(test, result.DataString, t.NewString(id))

	})

	test.Run("it should fail when I try to insert 2 times the same pk", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		example := &ExampleDataShared{
			DataString: t.NewString(id),
			DataInt:    t.NewInt64(1),
		}

		// When
		mgmTestData.Db.Create(context.Background(), t.NewString(id), example)
		err := mgmTestData.Db.Create(context.Background(), t.NewString(id), example)

		// Then
		assert.NotNil(test, err)
	})
}
