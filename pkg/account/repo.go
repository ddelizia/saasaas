package account

import (
	"context"

	"github.com/ddelizia/saasaas/pkg/t"
)

// AccountRepository is the interface to store the data (it will need to be implemented)
type Repository interface {
	Create(c context.Context, id t.String, data *Account) error
	Get(c context.Context, id t.String) (*Account, error)
	Update(c context.Context, id t.String, data *Account) error
	FindAll(c context.Context, startAt string, limit int64) (*AccountCursorList, error)
	Delete(c context.Context, id t.String) error
}
