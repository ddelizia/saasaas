package dynamodbsaas_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/ddelizia/saasaas/pkg/db/dynamodbsaas"
	"github.com/stretchr/testify/assert"
)

type AccountDataExample struct {
	Account string
	Status  string
}

func TestMarshalMap(t *testing.T) {
	repo := &dynamodbsaas.Db{
		EntityDescriptorMapping: map[reflect.Type]*dynamodbsaas.EntityDescriptor{
			reflect.TypeOf(&AccountDataExample{}): {
				DynamicFields: map[string]dynamodbsaas.DynamicAttrFunc{
					"key1": func(ctx context.Context, model interface{}) (interface{}, error) {
						example := model.(*AccountDataExample)
						return fmt.Sprintf("#ACCOUNT#%s", example.Account), nil
					},
					"key2": func(ctx context.Context, model interface{}) (interface{}, error) {
						example := model.(*AccountDataExample)
						return fmt.Sprintf("#STATUS#%s", example.Status), nil
					},
				},
			},
		},
	}

	t.Run("it should contain the generated key values", func(t *testing.T) {
		// Given
		exampleAccountData := &AccountDataExample{
			Account: "accountId",
			Status:  "Status",
		}

		// When
		result, err := repo.MarshalMap(context.Background(), nil, exampleAccountData)

		// Then
		assert.Contains(t, result["key1"].String(), "#ACCOUNT#accountId")
		assert.Contains(t, result["key2"].String(), "#STATUS#Status")
		assert.Nil(t, err)
	})

}
