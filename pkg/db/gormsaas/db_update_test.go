package gormsaas_test

import (
	"context"
	"testing"

	"github.com/ddelizia/saasaas/pkg/db/gormsaas"
	"github.com/ddelizia/saasaas/pkg/t"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateShared(test *testing.T) {
	gormTestData := SetupGorm(test)

	test.Run("it should fail when I try to update something that does not exists", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		example := &ExampleDataShared{
			GormSaasModel: gormsaas.GormSaasModel{ID: t.NewString(id)},
			DataString:    t.NewString(id),
			DataInt:       t.NewInt64(1),
		}

		// When
		err := gormTestData.Repo.Update(context.Background(), example.ID, example)

		// Then
		assert.NotNil(test, err)
	})

	test.Run("it should update the data correctly", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		dataToInsert := &ExampleDataShared{
			GormSaasModel: gormsaas.GormSaasModel{ID: t.NewString(id)},
			DataString:    t.NewString(id),
			DataInt:       t.NewInt64(1),
		}
		gormTestData.Gorm.Create(dataToInsert)

		// When
		example := &ExampleDataShared{
			GormSaasModel: gormsaas.GormSaasModel{ID: t.NewString(id)},
			DataString:    t.NewString(id),
			DataInt:       t.NewInt64(1),
		}
		err := gormTestData.Repo.Update(context.Background(), t.NewString(id), example)

		// Then
		assert.Nil(test, err)

		result := &ExampleDataShared{}
		err = gormTestData.Repo.Get(context.Background(), t.NewString(id), result)

		assert.Nil(test, err)
		assert.EqualValues(test, *result.DataString, *t.NewString("Hello"))
		assert.EqualValues(test, *result.DataInt, *t.NewInt64(2))

	})
}
