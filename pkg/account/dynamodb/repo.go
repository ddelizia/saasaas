package dynamodb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ddelizia/saasaas/pkg/account"
	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/db/dynamodbsaas"
	"github.com/ddelizia/saasaas/pkg/errors"
	"github.com/ddelizia/saasaas/pkg/t"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type accountModel struct {
	account.Account
	dynamodbsaas.DynamoSaasModel
}

type accountDb interface {
	db.Creator
	db.CursorLister
	db.Getter
	db.Updater
	db.Deleter
}

type baseAccountDb struct {
	dynamodbsaas.Db
}

type accountRepository struct {
	db accountDb
}

func NewRepository(db *dynamodb.DynamoDB, tableName string) account.Repository {
	const pkPrefix string = "#TENANT"
	const skPrefix string = "#ID"
	repoSharedDaf := map[string]dynamodbsaas.DynamicAttrFunc{}
	dynamodbsaas.StdDFWrapperSharedModel(pkPrefix, skPrefix, repoSharedDaf)

	pk := &dynamodbsaas.Index{
		PkField: "pk",
		SkField: "sk",
	}

	return &accountRepository{
		db: &baseAccountDb{
			Db: dynamodbsaas.Db{
				DynamoDb:  db,
				TableName: tableName,
				EntityDescriptorMapping: map[reflect.Type]*dynamodbsaas.EntityDescriptor{
					reflect.TypeOf(&accountModel{}): {
						DynamicFields: repoSharedDaf,
						IDBuilder:     dynamodbsaas.StdIDBuilderSharedModel(pkPrefix, skPrefix),
						PK:            pk,
					},
				},
			},
		},
	}
}

func (r *accountRepository) Get(c context.Context, accountID t.String) (*account.Account, error) {
	accountDynamo := &accountModel{}
	err := r.db.Get(c, accountID, accountDynamo)
	if err != nil {
		return nil, errors.NewTraceable(errors.NotFound, fmt.Sprintf("it account with id [%s] not found", accountID), err)
	}
	finalAccount := accountDynamo.Account
	return &finalAccount, nil
}

func (r *accountRepository) Create(c context.Context, id t.String, data *account.Account) error {
	err := r.db.Create(c, id, &accountModel{
		Account: *data,
	})
	if err != nil {
		return errors.NewTraceable(errors.CreateError, "error while creating the instance", err)
	}
	return nil
}

func (r *accountRepository) Update(c context.Context, id t.String, data *account.Account) error {
	err := r.db.Update(c, id, &accountModel{
		Account: *data,
	})
	if err != nil {
		return errors.NewTraceable(errors.UpdateError, "error while updating the instance", err)
	}
	return nil
}

func (r *accountRepository) FindAll(c context.Context, startAt string, limit int64) (*account.AccountCursorList, error) {
	var rr []*accountModel
	l, err := r.db.CursorList(c, startAt, limit, rr)
	if err != nil {
		return nil, errors.NewTraceable(errors.UpdateError, "error while executing query", err)
	}
	accounts := make([]*account.Account, 0)
	for _, r := range rr {
		accounts = append(accounts, &(r.Account))
	}

	return &account.AccountCursorList{
		LastKey: t.NewString(l.LastKey),
		Results: accounts,
	}, nil
}

func (r *accountRepository) Delete(c context.Context, id t.String) error {
	err := r.db.Delete(c, id, &accountModel{})
	if err != nil {
		return errors.NewTraceable(errors.DeleteError, "error while updating the instance", err)
	}
	return nil
}
