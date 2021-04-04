package gormsaas_test

import (
	"context"
	"testing"

	"github.com/ddelizia/saasaas/pkg/db/gormsaas"
	"github.com/ddelizia/saasaas/pkg/t"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type BaseStruct struct {
	DataInt    t.Int64
	DataString t.String
}

type ExtendedStruct struct {
	gormsaas.GormSaasModel
	BaseStruct
}

func TestCreateShared(test *testing.T) {
	gormTestData := SetupGorm(test)

	test.Run("it should work when I have a nested struct", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		example := &ExtendedStruct{
			BaseStruct: BaseStruct{
				DataString: t.NewString(id),
				DataInt:    t.NewInt64(1),
			},
		}
		gormTestData.Repo.GormDb.AutoMigrate(&ExtendedStruct{})

		// When
		err := gormTestData.Repo.Create(context.Background(), t.NewString(id), example)

		// Then
		assert.Nil(test, err)

		result := ExtendedStruct{}
		gormTestData.Gorm.Where("id = ?", id).First(&result)

		assert.Equal(test, result.DataInt, t.NewInt64(1))
		assert.Equal(test, result.DataString, t.NewString(id))
	})

	test.Run("it should return the data correctly", func(test *testing.T) {
		// Given
		id := uuid.New().String()
		example := &ExampleDataShared{
			DataString: t.NewString(id),
			DataInt:    t.NewInt64(1),
		}

		// When
		err := gormTestData.Repo.Create(context.Background(), t.NewString(id), example)

		// Then
		assert.Nil(test, err)

		result := ExampleDataShared{}
		gormTestData.Gorm.Where("id = ?", id).First(&result)

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
		gormTestData.Repo.Create(context.Background(), t.NewString(id), example)
		err := gormTestData.Repo.Create(context.Background(), t.NewString(id), example)

		// Then
		assert.NotNil(test, err)
	})

}
