package dynamodbsaas

import (
	"context"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/errors"
	"github.com/ddelizia/saasaas/pkg/t"
)

const (
	EntityMappingError errors.TraceableType = "EntityMappingError"
)

type DynamicAttrFunc func(ctx context.Context, model interface{}) (interface{}, error)
type IDBuilderFunc func(ctx context.Context, id t.String) map[string]*dynamodb.AttributeValue

type Index struct {
	PkField string
	SkField string
}

type EntityDescriptor struct {
	DynamicFields map[string]DynamicAttrFunc
	IDBuilder     IDBuilderFunc
	PK            *Index
}

type Db struct {
	DynamoDb                *dynamodb.DynamoDB
	TableName               string
	EntityDescriptorMapping map[reflect.Type]*EntityDescriptor
	db.Getter
	db.Creator
	db.Deleter
	db.Updater
	db.CursorLister
}

func (r *Db) GetEntityDescriptor(data interface{}) (*EntityDescriptor, error) {
	theType := reflect.TypeOf(data)

	if theType.Kind() == reflect.Ptr && theType.Elem().Kind() == reflect.Slice {
		theType = theType.Elem().Elem()
	}

	result := r.EntityDescriptorMapping[theType]
	if result == nil {
		return nil, errors.NewTraceable(EntityMappingError, fmt.Sprintf("not able to find entity [%v] ", reflect.TypeOf(data)), nil)
	}
	return result, nil
}

// Create instance entity data in Dynamodb table
func (r *Db) Create(ctx context.Context, id t.String, data interface{}) error {
	entityDescriptor, err := r.GetEntityDescriptor(data)
	if err != nil {
		return errors.NewTraceable(db.DataCreationError, "entity cannot be found", err)
	}

	marshalledData, err := r.MarshalMap(ctx, id, data)
	if err != nil {
		return errors.NewTraceable(db.DataCreationError, fmt.Sprintf("data [%v] cannot be unmarshalled", data), err)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item:      marshalledData,
	}

	input.ConditionExpression = aws.String(fmt.Sprintf("attribute_not_exists(%s)", entityDescriptor.PK.PkField))
	_, err = r.DynamoDb.PutItemWithContext(ctx, input)

	if err != nil {
		return errors.NewTraceable(db.DataCreationError, fmt.Sprintf("There was an error creating putting the item in the table [%v]", marshalledData), err)
	}

	return nil
}

// Get data from repository
func (r *Db) Get(ctx context.Context, id t.String, out interface{}) error {
	entityDescriptor, err := r.GetEntityDescriptor(out)
	if err != nil {
		return errors.NewTraceable(db.DataGetError, "entity cannot be found", err)
	}

	result, err := r.DynamoDb.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: &r.TableName,
		Key:       entityDescriptor.IDBuilder(ctx, id),
	})
	if err != nil {
		return errors.NewTraceable(db.DataGetError, "error getting data from dynamodb", err)
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, out)
	if err != nil {
		return errors.NewTraceable(db.DataGetError, fmt.Sprintf("error unmarshalling the data %v", result.Item), err)
	}
	return nil
}

// CursorList retrieves data from repository and returns as cursor
func (r *Db) CursorList(ctx context.Context, startAtKey string, limit int64, out interface{}) (*db.CursorListResult, error) {
	entityDescriptor, err := r.GetEntityDescriptor(out)
	if err != nil {
		return nil, errors.NewTraceable(db.DataListError, "entity cannot be found", err)
	}
	pkI, err := entityDescriptor.DynamicFields[entityDescriptor.PK.PkField](ctx, nil)
	if err != nil {

	}
	pkM, err := dynamodbattribute.Marshal(pkI)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.TableName),
		KeyConditionExpression: aws.String(fmt.Sprintf("%[1]s = :%[1]s", entityDescriptor.PK.PkField)),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":" + entityDescriptor.PK.PkField: pkM,
		},
	}

	return r.InternalCursorList(ctx, input, startAtKey, limit, out)
}

func (r *Db) InternalCursorList(ctx context.Context, input *dynamodb.QueryInput, startAtKey string, limit int64, out interface{}) (*db.CursorListResult, error) {

	//TODO: check the type of out argument
	t := reflect.TypeOf(out).Elem().Elem()

	//vp := v.Pointer()

	if limit != 0 {
		input.Limit = aws.Int64(limit)
	}

	if startAtKey != "" {
		d, err := DeserializeKey(startAtKey)
		if err != nil {
			return nil, errors.NewTraceable(db.DataListError, fmt.Sprintf("cannot deserialize key [%s]", startAtKey), err)
		}
		input.ExclusiveStartKey = d
	}

	q, err := r.DynamoDb.QueryWithContext(ctx, input)

	if err != nil {
		return nil, errors.NewTraceable(db.DataListError, "error executing the query", err)
	}

	lt := reflect.SliceOf(t)
	lv := reflect.New(lt).Elem()
	for _, k := range q.Items {
		v := reflect.New(t)

		concreteObject := v.Interface()
		err = dynamodbattribute.UnmarshalMap(k, concreteObject)
		if err != nil {
			return nil, errors.NewTraceable(db.DataListError, fmt.Sprintf("not able to unmarshal %v", k), err)
		}
		v = reflect.ValueOf(concreteObject).Elem()
		lv = reflect.Append(lv, v)
	}

	outV := reflect.ValueOf(out).Elem()
	outV.Set(lv)

	lk, err := SerializeKey(q.LastEvaluatedKey)

	if err != nil {
		return nil, errors.NewTraceable(db.DataListError, "not able to serialize last key", err)
	}
	return &db.CursorListResult{
		LastKey: lk,
	}, nil
}

// Update data into repository
func (r *Db) Update(ctx context.Context, id t.String, data interface{}) error {
	entityDescriptor, err := r.GetEntityDescriptor(data)
	if err != nil {
		return errors.NewTraceable(db.DataUpdateError, "entity cannot be found", err)
	}

	marshalledData, err := r.MarshalMap(ctx, id, data)
	if err != nil {
		return errors.NewTraceable(db.DataUpdateError, fmt.Sprintf("data [%v] cannot be unmarshalled", data), err)
	}
	in := &dynamodb.PutItemInput{
		TableName: &r.TableName,
		Item:      marshalledData,
	}

	in.ConditionExpression = aws.String(fmt.Sprintf("attribute_exists(%s)", entityDescriptor.PK.PkField))

	_, err = r.DynamoDb.PutItemWithContext(ctx, in)

	if err != nil {
		return errors.NewTraceable(db.DataUpdateError, fmt.Sprintf("there was an error updating the data [%v]", marshalledData), err)
	}

	return nil
}

// Delete data from repository
func (r *Db) Delete(ctx context.Context, id t.String, out interface{}) error {
	entityDescriptor, err := r.GetEntityDescriptor(out)
	if err != nil {
		return errors.NewTraceable(db.DataDeleteError, "entity cannot be found", err)
	}

	in := &dynamodb.DeleteItemInput{
		TableName: &r.TableName,
		Key:       entityDescriptor.IDBuilder(ctx, id),
	}

	_, err = r.DynamoDb.DeleteItemWithContext(ctx, in)

	if err != nil {
		return errors.NewTraceable(db.DataDeleteError, fmt.Sprintf("there was an error while deleting item with id [%v]", id), err)
	}

	return nil
}
