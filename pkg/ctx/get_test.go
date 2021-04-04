package ctx_test

import (
	"context"
	"testing"

	"github.com/ddelizia/saasaas/pkg/ctx"
	"github.com/stretchr/testify/assert"
)

type Example struct {
	Data string
}

func TestGetFromContext(t *testing.T) {
	t.Run("it should should get the data from the context", func(t *testing.T) {
		// Given
		initialData := &Example{
			Data: "hello",
		}
		var key ctx.ContextKey = "data"
		aContext := context.WithValue(context.Background(), key, initialData)

		// When
		result, err := ctx.GetFromContext(aContext, "data")

		// Then
		assert.Equal(t, result, initialData)
		assert.Nil(t, err)
	})

	t.Run("it should should return an error when value is empty", func(t *testing.T) {
		// Given
		aContext := context.Background()

		// When
		var key ctx.ContextKey = "data"
		result, err := ctx.GetFromContext(aContext, key)

		// Then
		assert.Nil(t, result)
		assert.NotNil(t, err)
	})
}
