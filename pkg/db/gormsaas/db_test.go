package gormsaas_test

import (
	"testing"

	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/db/gormsaas"
	"github.com/ddelizia/saasaas/pkg/t"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormTestSetup struct {
	Repo *gormsaas.Db
	Gorm *gorm.DB
}

type ExampleDataShared struct {
	gormsaas.GormSaasModel
	DataInt    t.Int64
	DataString t.String
}

type ExampleDataTenant struct {
	gormsaas.GormSaasModel
	db.AccountAware
	DataInt    t.Int64
	DataString t.String
}

type ExampleDataUser struct {
	gormsaas.GormSaasModel
	db.UserAware
	DataInt    t.Int64
	DataString t.String
}

func SetupGorm(t *testing.T) *GormTestSetup {
	gormDb, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	gormDb.AutoMigrate(&ExampleDataShared{})
	gormDb.AutoMigrate(&ExampleDataTenant{})
	gormDb.AutoMigrate(&ExampleDataUser{})

	repo := &gormsaas.Db{
		GormDb: gormDb,
	}

	return &GormTestSetup{
		Repo: repo,
		Gorm: gormDb,
	}
}
