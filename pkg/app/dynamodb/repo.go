package dynamodb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ddelizia/saasaas/pkg/app"
	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/db/dynamodbsaas"
	"github.com/ddelizia/saasaas/pkg/errors"
)

const (
	lsi1 string = "lsi1"
)

type appModel struct {
	app.App
	dynamodbsaas.DynamoSaasModel
}

type appDb interface {
	db.Creator
	db.CursorLister
	db.Getter
	db.Updater
	db.Deleter
	CursorListByTenant(ctx context.Context, tenantID string, startAtKey string, limit int64, out interface{}) (*db.CursorListResult, error)
}

type baseAppDb struct {
	dynamodbsaas.Db
}

func (r *baseAppDb) CursorListByTenant(ctx context.Context, tenantID string, startAtKey string, limit int64, out interface{}) (*db.CursorListResult, error) {
	entityDescriptor, err := r.GetEntityDescriptor(out)
	if err != nil {
		
	}
	pkI, err := entityDescriptor.DynamicFields[r.PK.PkField](ctx, nil)
	if err != nil {

	}
	pkM, err := dynamodbattribute.Marshal(pkI)

	lsi1I, err := entityDescriptor.DynamicFields[r.PK.PkField](ctx, &appModel{
		App: app.App{
			TenantID: tenantID,
		},
	})
	if err != nil {

	}
	lsi1M, err := dynamodbattribute.Marshal(lsi1I)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String(fmt.Sprintf("%[1]s = :%[1]s and %[2]s = %[2]s", r.PK.PkField, lsi1)),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":" + r.PK.PkField: pkM,
			":" + lsi1:         lsi1M,
		},
	}

	return r.InternalCursorList(ctx, input, startAtKey, limit, out)
}

type appRepository struct {
	db appDb
}

func NewAppDynamoRepository(db *dynamodb.DynamoDB, tableName string) app.Repository {
	const pkPrefix string = "#APP"
	const skPrefix string = "#ID"
	repoSharedDaf := map[string]dynamodbsaas.DynamicAttrFunc{
		lsi1: func(ctx context.Context, model interface{}) (interface{}, error) {
			app := model.(*appModel)
			return fmt.Sprintf("#TENANT#%s", app.TenantID), nil
		},
	}
	dynamodbsaas.StdDFWrapperSharedModel(pkPrefix, skPrefix, repoSharedDaf)

	return &appRepository{
		db: &baseAppDb{
			Db: dynamodbsaas.Db{
				DynamoDb: db,
				PK: &dynamodbsaas.Index{
					PkField: "pk",
					SkField: "sk",
				},
				TableName: tableName,
				EntityDescriptorMapping: map[reflect.Type]*dynamodbsaas.EntityDescriptor{
					reflect.TypeOf(&appModel{}): {
						DynamicFields: repoSharedDaf,
						IDBuilder:     dynamodbsaas.StdIDBuilderSharedModel(pkPrefix, skPrefix),
					},
				},
			},
		},
	}
}

func (r *appRepository) Create(c context.Context, id string, data *app.App) error {
	err := r.db.Create(c, id, &appModel{
		App: *data,
	})
	if err != nil {
		return errors.NewTraceable(errors.CreateError, "error while creating the instance", err)
	}
	return nil
}

func (r *appRepository) FindAll(c context.Context, startAt string, limit int64) (*app.AppCursorList, error) {
	var rr []*appModel
	l, err := r.db.CursorList(c, startAt, limit, rr)
	if err != nil {
		return nil, errors.NewTraceable(errors.UpdateError, "error while executing query", err)
	}
	tenants := make([]*app.App, 0)
	for _, r := range rr {
		tenants = append(tenants, &(r.App))
	}

	return &app.AppCursorList{
		LastKey: l.LastKey,
		Results: tenants,
	}, nil
}

func (r *appRepository) FindByTenant(c context.Context, tenantID string, startAt string, limit int64) (*app.AppCursorList, error) {
	var rr []*appModel
	l, err := r.db.CursorListByTenant(c, startAt, tenantID, limit, rr)
	if err != nil {
		return nil, errors.NewTraceable(errors.UpdateError, "error while executing query", err)
	}
	tenants := make([]*app.App, 0)
	for _, r := range rr {
		tenants = append(tenants, &(r.App))
	}

	return &app.AppCursorList{
		LastKey: l.LastKey,
		Results: tenants,
	}, nil
}

func (r *appRepository) Delete(c context.Context, id string) error {
	err := r.db.Delete(c, id, &appModel{})
	if err != nil {
		return errors.NewTraceable(errors.DeleteError, "error while updating the instance", err)
	}
	return nil
}

func (r *appRepository) Update(c context.Context, id string, data *app.App) error {
	err := r.db.Update(c, id, &appModel{
		App: *data,
	})
	if err != nil {
		return errors.NewTraceable(errors.UpdateError, "error while updating the instance", err)
	}
	return nil
}

func (r *appRepository) Get(c context.Context, tenantID string) (*app.App, error) {
	tenantDynamo := &appModel{}
	err := r.db.Get(c, tenantID, tenantDynamo)
	if err != nil {
		return nil, errors.NewTraceable(errors.NotFound, fmt.Sprintf("it tenant with id [%s] not found", tenantID), err)
	}
	finalTenant := tenantDynamo.App
	return &finalTenant, nil
}
