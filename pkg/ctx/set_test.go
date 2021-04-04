package ctx_test

import (
	"context"
	"testing"

	"github.com/ddelizia/saasaas/pkg/ctx"
	"github.com/stretchr/testify/assert"
)

func TestSetInContext(t *testing.T) {
	t.Run("it should should set the data in the context", func(t *testing.T) {
		// Given
		newContext := context.Background()

		// When
		var key ctx.ContextKey = "data"
		ctxResult := ctx.SetInContext(newContext, key, "hello")

		// Then
		assert.Equal(t, ctxResult.Value(key), "hello")
	})
}
