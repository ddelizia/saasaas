package ctx

import (
	"context"
	"fmt"

	"github.com/ddelizia/saasaas/pkg/errors"
)

// GetFromContext set tenant data in context
func GetFromContext(ctx context.Context, key ContextKey) (interface{}, error) {
	result := ctx.Value(key)
	if result == nil {
		return nil, errors.NewTraceable(errors.NotFound, fmt.Sprintf("not able to find the key [%s] in context", key), nil)
	}
	return result, nil
}
