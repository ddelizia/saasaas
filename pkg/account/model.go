package account

import "github.com/ddelizia/saasaas/pkg/t"

// Account represents the account information
type Account struct {
	ID   t.String
	Name t.String
}

type AccountCursorList struct {
	LastKey t.String
	Results []*Account
}
