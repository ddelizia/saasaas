package dynamodbsaas

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ddelizia/saasaas/pkg/ctx"
	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/errors"
	"github.com/ddelizia/saasaas/pkg/t"
)

const MarshalingError errors.TraceableType = "MarshalingError"

// MarshalMap execute wrapper on the object
func (r *Db) MarshalMap(c context.Context, id t.String, data interface{}) (map[string]*dynamodb.AttributeValue, error) {

	marshalledData, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return nil, errors.NewTraceable(MarshalingError, fmt.Sprintf("data [%v] cannot be marshalled", data), err)
	}

	entityDescriptor, err := r.GetEntityDescriptor(data)
	if err != nil {
		return nil, errors.NewTraceable(MarshalingError, "entity cannot be found", err)
	}

	accountFieldValue := reflect.ValueOf(data).Elem().FieldByName(db.ModelFieldAccountID)
	if accountFieldValue.IsValid() {
		account, _ := ctx.GetFromContext(c, ctx.AccountIDContextField)
		accountFieldValue.SetString(account.(string))
	}

	userFieldValue := reflect.ValueOf(data).Elem().FieldByName(db.ModelFieldUserID)
	if userFieldValue.IsValid() {
		user, _ := ctx.GetFromContext(c, ctx.UserIDContextField)
		accountFieldValue.SetString(user.(string))
	}

	appFieldValue := reflect.ValueOf(data).Elem().FieldByName(db.ModelFieldAppID)
	if appFieldValue.IsValid() {
		app, _ := ctx.GetFromContext(c, ctx.ProjectIDContextField)
		appFieldValue.SetString(app.(string))
	}

	for k, da := range entityDescriptor.DynamicFields {
		evalKey, err := da(c, data)
		if err != nil {
			return nil, errors.NewTraceable(MarshalingError, fmt.Sprintf("Error evaluating function for key [%s]", k), err)
		}

		currentMarshalledData, err := dynamodbattribute.Marshal(evalKey)
		if err != nil {
			return nil, errors.NewTraceable(MarshalingError, fmt.Sprintf("data in evaluted key [%s] and value [%v] cannot be marshalled", k, evalKey), err)
		}
		marshalledData[k] = currentMarshalledData
	}

	return marshalledData, nil
}
