package stripe

import "github.com/ddelizia/saasaas/pkg/t"

type Customer struct {
	ID        t.String
	AccountID t.String
	PlanID    t.String
}
