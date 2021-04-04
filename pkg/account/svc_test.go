package account_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ddelizia/saasaas/pkg/account"
	"github.com/ddelizia/saasaas/pkg/t"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//******************
//*** INIT MOCKS ***
//******************
type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) Get(ctx context.Context, accountID t.String) (*account.Account, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*account.Account), args.Error(1)
}

func (m *RepositoryMock) Create(ctx context.Context, id t.String, data *account.Account) error {
	args := m.Called(ctx, id, data)
	return args.Error(0)
}

func (m *RepositoryMock) Delete(ctx context.Context, id t.String) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *RepositoryMock) Update(ctx context.Context, id t.String, data *account.Account) error {
	args := m.Called(ctx, id, data)
	return args.Error(0)
}

func (m *RepositoryMock) FindAll(c context.Context, startAt string, limit int64) (*account.AccountCursorList, error) {
	args := m.Called(c, startAt, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*account.AccountCursorList), args.Error(1)
}

type testService struct {
	service  account.Service
	accounts *RepositoryMock
}

func newTestService() *testService {
	accounts := &RepositoryMock{}
	service := account.New(accounts)

	return &testService{
		service:  service,
		accounts: accounts,
	}
}

//******************
//***** TESTS ******
//******************

func TestUpdate(tt *testing.T) {
	exampleData := &account.Account{
		ID:   t.NewString("id"),
		Name: t.NewString("name"),
	}

	tt.Run("it should call the repository", func(tt *testing.T) {
		// Given
		testSvc := newTestService()
		testSvc.accounts.On("Update", mock.Anything, t.NewString("id"), exampleData).Return(nil)

		// When
		err := testSvc.service.Update(context.Background(), t.NewString("id"), exampleData)

		// Then
		assert.Nil(tt, err)
		testSvc.accounts.AssertNumberOfCalls(tt, "Update", 1)
	})

}

func TestGet(tt *testing.T) {
	exampleData := &account.Account{
		ID:   t.NewString("id"),
		Name: t.NewString("name"),
	}

	tt.Run("it should return the data correctly", func(tt *testing.T) {
		// Given
		testSvc := newTestService()
		testSvc.accounts.On("Get", mock.Anything, t.NewString("id")).Return(exampleData, nil)

		// When
		result, err := testSvc.service.Get(context.Background(), t.NewString("id"))

		// Then
		assert.Equal(tt, result, exampleData)
		assert.Nil(tt, err)
	})

	tt.Run("it should return error when data cannot be found in the repository", func(tt *testing.T) {
		// Given
		testSvc := newTestService()
		testSvc.accounts.On("Get", mock.Anything, t.NewString("id")).Return(nil, fmt.Errorf("data not found"))

		// When
		result, err := testSvc.service.Get(context.Background(), t.NewString("id"))

		// Then
		assert.Nil(tt, result)
		assert.NotNil(tt, err)
	})
}

func TestCreate(tt *testing.T) {

	exampleData := &account.Account{
		ID:   t.NewString("id"),
		Name: t.NewString("name"),
	}

	tt.Run("it should create correctly the account", func(ttt *testing.T) {
		// Given
		testSvc := newTestService()
		testSvc.accounts.On("Create", mock.Anything, t.NewString("id"), exampleData).Return(nil)

		// When
		id, err := testSvc.service.Create(context.Background(), exampleData)

		// Then
		assert.Equal(tt, id, t.NewString("id"))
		assert.Nil(tt, err)
	})

	tt.Run("it should fail when repository creation fails", func(tt *testing.T) {
		// Given
		testSvc := newTestService()
		exampleData := &account.Account{
			ID:   t.NewString("id2"),
			Name: t.NewString("name"),
		}
		testSvc.accounts.On("Create", mock.Anything, t.NewString("id2"), exampleData).Return(fmt.Errorf("throw this error"))

		// When
		id, err := testSvc.service.Create(context.Background(), exampleData)

		// Then
		assert.Nil(tt, id)
		assert.NotNil(tt, err)
	})

	tt.Run("it should call the repository only once", func(tt *testing.T) {
		// Given
		testSvc := newTestService()
		exampleData := &account.Account{
			ID:   t.NewString("id3"),
			Name: t.NewString("name"),
		}
		testSvc.accounts.On("Create", mock.Anything, t.NewString("id3"), exampleData).Return(nil)

		// When
		testSvc.service.Create(context.Background(), exampleData)

		// Then
		testSvc.accounts.AssertNumberOfCalls(tt, "Create", 1)
	})
}

func TestDelete(tt *testing.T) {

	tt.Run("it should get an error", func(tt *testing.T) {
		// Given
		testSvc := newTestService()
		testSvc.accounts.On("Delete", mock.Anything, t.NewString("id")).Return(errors.New("Example error"))

		// When
		err := testSvc.service.Delete(context.Background(), t.NewString("id"))

		// Then
		testSvc.accounts.AssertNumberOfCalls(tt, "Delete", 1)
		assert.NotNil(tt, err)
	})

	tt.Run("it should call delete", func(tt *testing.T) {
		// Given
		testSvc := newTestService()
		testSvc.accounts.On("Delete", mock.Anything, t.NewString("id")).Return(nil)

		// When
		testSvc.service.Delete(context.Background(), t.NewString("id"))

		// Then
		testSvc.accounts.AssertNumberOfCalls(tt, "Delete", 1)
	})
}
