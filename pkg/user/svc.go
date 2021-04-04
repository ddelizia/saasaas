package user

import "context"

type Service interface {
	Create(c context.Context, data *User) (string, error)
	Get(c context.Context, tenantID string) (*User, error)
	Update(c context.Context, tenantID string, data *User) error
	Delete(c context.Context, tenantID string) error
	
}
