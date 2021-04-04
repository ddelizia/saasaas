package account

import (
	"context"
	"fmt"

	"github.com/ddelizia/saasaas/pkg/errors"
	"github.com/ddelizia/saasaas/pkg/t"
	"github.com/google/uuid"
)

const (
	// UpdateError update error
	UpdateError errors.TraceableType = "AccountUpdateError"
	// GetError when account cannot be retrieved
	GetError errors.TraceableType = "AccountGetError"
	// CreationError is thrown when creation of the account is not satisfactory
	CreationError errors.TraceableType = "AccountCreationError"
	// DeleteError when account cannot be retrieved
	DeleteError errors.TraceableType = "AccountDeleteError"
	// ListError when account cannot be found
	ListError errors.TraceableType = "AccountDeleteError"
)

// Service is the inteface that exposes the services
type Service interface {
	Create(c context.Context, data *Account) (t.String, error)
	Get(c context.Context, accountID t.String) (*Account, error)
	Update(c context.Context, accountID t.String, data *Account) error
	Delete(c context.Context, accountID t.String) error
	List(c context.Context, startAt string, limit int64) (*AccountCursorList, error)
}

// Service implementation that abstract the repository
type service struct {
	accounts Repository
}

// New creates a new instance of the service
func New(accounts Repository) Service {
	return &service{
		accounts: accounts,
	}
}

// Get account information
func (s *service) Get(ctx context.Context, accountID t.String) (*Account, error) {
	result, err := s.accounts.Get(ctx, accountID)
	if err != nil {
		return nil, errors.NewTraceable(GetError, "it was not able to find the account in the repository", err)
	}
	return result, nil
}

// Update account information
func (s *service) Update(ctx context.Context, accountID t.String, data *Account) error {
	err := s.accounts.Update(ctx, accountID, data)
	if err != nil {
		return errors.NewTraceable(UpdateError, "it was not able to find the account in the repository", err)
	}
	return nil
}

// Create account information
func (s *service) Create(ctx context.Context, data *Account) (t.String, error) {
	id := data.ID
	if id == nil {
		uuidValue, err := uuid.NewUUID()
		if err != nil {
			return nil, errors.NewTraceable(CreationError, "not able to create id via UUID for account", err)
		}
		id = t.NewString(uuidValue.String())
	}

	err := s.accounts.Create(ctx, data.ID, data)
	if err != nil {
		return nil, errors.NewTraceable(CreationError, "error during the creation of the repository", err)
	}
	return id, nil
}

func (s *service) Delete(c context.Context, accountID t.String) error {
	err := s.accounts.Delete(c, accountID)

	if err != nil {
		return errors.NewTraceable(DeleteError, fmt.Sprintf("it was not able to delete the account with id [%s]", *accountID), err)
	}

	return nil
}

func (s *service) List(c context.Context, startAt string, limit int64) (*AccountCursorList, error) {
	result, err := s.accounts.FindAll(c, startAt, limit)
	if err != nil {
		return nil, errors.NewTraceable(DeleteError, "it was not able to find the data", err)
	}

	return result, err
}
