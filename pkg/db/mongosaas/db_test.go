package mongosaas_test

import (
	"context"
	"testing"

	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/db/mongosaas"
	"github.com/ddelizia/saasaas/pkg/t"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoTestSetup struct {
	Db     *mongosaas.Db
	client *mongo.Client
}

type ExampleDataShared struct {
	mongosaas.MongoSaasModel
	DataInt    t.Int64
	DataString t.String
}

type ExampleDataTenant struct {
	mongosaas.MongoSaasModel
	db.AccountAware
	DataInt    t.Int64
	DataString t.String
}

type ExampleDataUser struct {
	mongosaas.MongoSaasModel
	db.UserAware
	DataInt    t.Int64
	DataString t.String
}

func SetupMongo(t *testing.T) *MongoTestSetup {
	ctx := context.Background()
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	//defer func() {
	//	if err := client.Disconnect(ctx); err != nil {
	//		panic(err)
	//	}
	//}()
	return &MongoTestSetup{
		Db: &mongosaas.Db{
			Client: client,
			DBName: "testing",
		},
		client: client,
	}
}
