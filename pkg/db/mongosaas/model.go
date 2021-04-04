package mongosaas

import (
	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/t"
)

type MongoSaasModel struct {
	db.SaasModel
	ID t.String
}
