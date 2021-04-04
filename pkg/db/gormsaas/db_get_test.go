package gormsaas_test

import (
	"context"
	"testing"

	"github.com/ddelizia/saasaas/pkg/db/gormsaas"
	"github.com/ddelizia/saasaas/pkg/t"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetShared(test *testing.T) {
	gormTestData := SetupGorm(test)

	test.Run("it should return correctly the data", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		dataToInsert := &ExampleDataShared{
			GormSaasModel: gormsaas.GormSaasModel{ID: t.NewString(id)},
			DataString:    t.NewString(id),
			DataInt:       t.NewInt64(1),
		}
		gormTestData.Gorm.Create(dataToInsert)

		// When
		result := &ExampleDataShared{}
		err := gormTestData.Repo.Get(context.Background(), t.NewString(id), result)

		// Then
		assert.Nil(test, err)
		assert.Equal(test, result.DataInt, t.NewInt64(1))
		assert.Equal(test, result.DataString, t.NewString(id))

	})

}
