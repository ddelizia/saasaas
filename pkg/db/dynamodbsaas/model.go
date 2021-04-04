package dynamodbsaas

import (
	"github.com/ddelizia/saasaas/pkg/db"
)

type DynamoSaasModel struct {
	db.SaasModel
	UpdatedAt int64
	CreatedAt int64
}
