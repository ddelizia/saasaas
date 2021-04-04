package app

import "github.com/ddelizia/saasaas/pkg/t"

type App struct {
	ID       t.String
	TenantID t.String
}

type AppCursorList struct {
	LastKey t.String
	Results []*App
}
