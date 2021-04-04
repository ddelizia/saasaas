package app

import "context"

// Repository is the interface to store the data (it will need to be implemented)
type Repository interface {
	Create(c context.Context, id string, data *App) error
	Get(c context.Context, id string) (*App, error)
	Update(c context.Context, id string, data *App) error
	Delete(c context.Context, id string) error
	FindAll(c context.Context, startAt string, limit int64) (*AppCursorList, error)
	FindByTenant(c context.Context, tenantID string, startAt string, limit int64) (*AppCursorList, error)
}
