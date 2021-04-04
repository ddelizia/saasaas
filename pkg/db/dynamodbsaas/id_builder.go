package dynamodbsaas

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ddelizia/saasaas/pkg/ctx"
	"github.com/ddelizia/saasaas/pkg/t"
)

// StdIDBuilderSharedModel build the standard ids
func StdIDBuilderSharedModel(pkPrefix string, skPrefix string) IDBuilderFunc {
	f := func(ctx context.Context, id t.String) map[string]*dynamodb.AttributeValue {
		pk, _ := dynamodbattribute.Marshal(pkPrefix)
		sk, _ := dynamodbattribute.Marshal(fmt.Sprintf("%s#%v", skPrefix, id))
		return map[string]*dynamodb.AttributeValue{
			"pk": pk,
			"sk": sk,
		}
	}
	return f
}

// StdIDBuilderAccountModel build the standard ids
func StdIDBuilderAccountModel(skPrefix string) IDBuilderFunc {
	f := func(c context.Context, id t.String) map[string]*dynamodb.AttributeValue {
		account, _ := ctx.GetFromContext(c, ctx.AccountIDContextField)
		pk, _ := dynamodbattribute.Marshal(fmt.Sprintf("#ACCOUNT#%s", account))
		sk, _ := dynamodbattribute.Marshal(fmt.Sprintf("%s#%v", skPrefix, id))
		return map[string]*dynamodb.AttributeValue{
			"pk": pk,
			"sk": sk,
		}
	}
	return f
}

// StdIDBuilderUserModel build the standard ids
func StdIDBuilderUserModel(skPrefix string) IDBuilderFunc {
	f := func(c context.Context, id t.String) map[string]*dynamodb.AttributeValue {
		user, _ := ctx.GetFromContext(c, ctx.UserIDContextField)
		pk, _ := dynamodbattribute.Marshal(fmt.Sprintf("#USER#%s", user))
		sk, _ := dynamodbattribute.Marshal(fmt.Sprintf("%s#%v", skPrefix, id))
		return map[string]*dynamodb.AttributeValue{
			"pk": pk,
			"sk": sk,
		}
	}
	return f
}

// StdIDBuilderAccountUserModel build the standard ids
func StdIDBuilderAccountUserModel(skPrefix string) IDBuilderFunc {
	f := func(c context.Context, id t.String) map[string]*dynamodb.AttributeValue {
		account, _ := ctx.GetFromContext(c, ctx.AccountIDContextField)
		user, _ := ctx.GetFromContext(c, ctx.UserIDContextField)
		pk, _ := dynamodbattribute.Marshal(fmt.Sprintf("#ACCOUNT#%s#USER#%s", account, user))
		sk, _ := dynamodbattribute.Marshal(fmt.Sprintf("%s#%v", skPrefix, id))
		return map[string]*dynamodb.AttributeValue{
			"pk": pk,
			"sk": sk,
		}
	}
	return f
}
