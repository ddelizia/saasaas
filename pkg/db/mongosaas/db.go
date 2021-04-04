package mongosaas

import (
	"context"
	"reflect"

	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/t"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Db struct {
	Client *mongo.Client
	DBName string
	db.Getter
	db.Creator
	db.Deleter
	db.Updater
	db.CursorLister
}

/*
func EnsureIndex(ctx context.Context, c *mongo.Collection, keys []string, opts *mongo.IndexOptionsBuilder) error {
	idxs := c.Indexes()

	ks := bson.D{}
	for _, k := range keys {
		// todo - add support for sorting index.
		ks.Append({})
	}
	idm := mongo.IndexModel{
		Keys:    ks,
		Options: opts.Build(),
	}

	v := idm.Options.Lookup("name")
	if v == nil {
		return errors.New("must provide a key name for index")
	}
	expectedName := v.StringValue()

	cur, err := idxs.List(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to list indexes")
	}

	found := false
	for cur.Next(ctx) {
		d := bson.NewDocument()

		if err := cur.Decode(d); err != nil {
			return errors.Wrap(err, "unable to decode bson index document")
		}

		v := d.Lookup("name")
		if v != nil && v.StringValue() == expectedName {
			found = true
			break
		}
	}

	if found {
		return nil
	}

	_, err = idxs.CreateOne(ctx, idm)
	return err
}
*/
func (r *Db) Init() {
	//EnsureIndex()
}

func (r *Db) Coll(data interface{}) *mongo.Collection {
	return r.Client.Database(r.DBName).Collection(reflect.TypeOf(data).Name())
}

func (r *Db) Create(ctx context.Context, id t.String, data interface{}) error {
	reflect.ValueOf(data).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
	_, err := r.Coll(data).InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	return err
}

func (r *Db) Get(ctx context.Context, id t.String, out interface{}) error {
	return nil
}

func (r *Db) Delete(ctx context.Context, id t.String, out interface{}) error {
	return nil
}

func (r *Db) Update(ctx context.Context, id t.String, data interface{}) error {
	return nil
}

func (r *Db) CursorList(ctx context.Context, startAtKey string, limit int64, out interface{}) (*db.CursorListResult, error) {
	return nil, nil
}
