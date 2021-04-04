package gormsaas

import (
	"github.com/ddelizia/saasaas/pkg/db"
	"github.com/ddelizia/saasaas/pkg/t"
	"gorm.io/gorm"
)

type GormSaasModel struct {
	gorm.Model
	db.SaasModel
	ID t.String
}
