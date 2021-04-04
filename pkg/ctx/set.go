package ctx

import (
	"context"
)

// SetInContext set info in context
func SetInContext(ctx context.Context, key ContextKey, data interface{}) context.Context {
	return context.WithValue(ctx, key, data)
}
