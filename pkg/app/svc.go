package app

import (
	"context"
	"fmt"

	"github.com/ddelizia/saasaas/pkg/errors"
)

const (
	// ListError when app cannot be listed
	ListError errors.TraceableType = "AppListError"
)

// Service for the app module
type Service interface {
	//Create(c context.Context, data *App) (string, error)
	//Get(c context.Context, appID string) (*App, error)
	//Update(c context.Context, appID string, data *App) error
	//Delete(c context.Context, appID string) error
	List(c context.Context, startAt string, limit int64) (*AppCursorList, error)
	GetAppsForTenant(c context.Context, tenantID string, startAt string, limit int64) (*AppCursorList, error)
}

type service struct {
	apps Repository
}

func New(apps Repository) Service {
	return &service{
		apps: apps,
	}
}

func (s *service) List(c context.Context, startAt string, limit int64) (*AppCursorList, error) {
	result, err := s.apps.FindAll(c, startAt, limit)
	if err != nil {
		return nil, errors.NewTraceable(ListError, "find has thrown an error", err)
	}

	return result, err
}

func (s *service) GetAppsForTenant(c context.Context, tenantID string, startAt string, limit int64) (*AppCursorList, error) {
	result, err := s.apps.FindAll(c, startAt, limit)
	if err != nil {
		return nil, errors.NewTraceable(ListError, fmt.Sprintf("it was not able to find data for tenant [%s]", tenantID), err)
	}

	return result, err
}
