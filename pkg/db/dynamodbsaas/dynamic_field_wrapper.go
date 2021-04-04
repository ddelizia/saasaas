package dynamodbsaas

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ddelizia/saasaas/pkg/ctx"
)

// StdDFWrapperSharedModel is a wrapper that adds "Standard" dynamic fields
func StdDFWrapperSharedModel(pkPrefix string, skPrefix string, da map[string]DynamicAttrFunc) {
	da["pk"] = func(c context.Context, model interface{}) (interface{}, error) {
		return pkPrefix, nil
	}

	da["sk"] = func(c context.Context, model interface{}) (interface{}, error) {
		return fmt.Sprintf("%s#%s", skPrefix, reflect.ValueOf(model).Elem().FieldByName("ID").String()), nil
	}
}

// StdDFWrapperAccountModel is a wrapper that adds "Standard" dynamic fields
func StdDFWrapperAccountModel(skPrefix string, da map[string]DynamicAttrFunc) {
	da["pk"] = func(c context.Context, model interface{}) (interface{}, error) {
		account, _ := ctx.GetFromContext(c, ctx.AccountIDContextField)
		return fmt.Sprintf("#ACCOUNT#%s", account), nil
	}

	da["sk"] = func(c context.Context, model interface{}) (interface{}, error) {
		return fmt.Sprintf("%s#%s", skPrefix, reflect.ValueOf(model).Elem().FieldByName("ID").String()), nil
	}
}

// StdDFWrapperUserModel is a wrapper that adds "Standard" dynamic fields
func StdDFWrapperUserModel(pkPrefix string, skPrefix string, da map[string]DynamicAttrFunc) {
	da["pk"] = func(c context.Context, model interface{}) (interface{}, error) {
		user, _ := ctx.GetFromContext(c, ctx.UserIDContextField)
		return fmt.Sprintf("#USER#%s", user), nil
	}

	da["sk"] = func(c context.Context, model interface{}) (interface{}, error) {
		return fmt.Sprintf("%s#%s", skPrefix, reflect.ValueOf(model).Elem().FieldByName("ID").String()), nil
	}
}

// StdDFWrapperAccountUserModel is a wrapper that adds "Standard" dynamic fields
func StdDFWrapperAccountUserModel(pkPrefix string, skPrefix string, da map[string]DynamicAttrFunc) {
	da["pk"] = func(c context.Context, model interface{}) (interface{}, error) {
		account, _ := ctx.GetFromContext(c, ctx.AccountIDContextField)
		user, _ := ctx.GetFromContext(c, ctx.UserIDContextField)
		return fmt.Sprintf("#ACCOUNT#%s#USER#%s", account, user), nil
	}

	da["sk"] = func(c context.Context, model interface{}) (interface{}, error) {
		return fmt.Sprintf("%s#%s", skPrefix, reflect.ValueOf(model).Elem().FieldByName("ID").String()), nil
	}
}
