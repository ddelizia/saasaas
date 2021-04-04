package db

import (
	"context"

	"github.com/ddelizia/saasaas/pkg/t"
)

type Getter interface {
	Get(ctx context.Context, id t.String, out interface{}) error
}

type Creator interface {
	Create(ctx context.Context, id t.String, data interface{}) error
}

type Updater interface {
	Update(ctx context.Context, id t.String, data interface{}) error
}

type CursorListResult struct {
	LastKey string
}

type CursorLister interface {
	CursorList(ctx context.Context, startAtKey string, limit int64, out interface{}) (*CursorListResult, error)
}

type PaginatedListResult struct {
	CurrentPage int64
	PageSize    int64
}

type PaginatedLister interface {
	PaginatedList(ctx context.Context, pageSize int64, currentPage int64, out []interface{}) (*PaginatedListResult, error)
}

type Deleter interface {
	Delete(ctx context.Context, id t.String, out interface{}) error
}
