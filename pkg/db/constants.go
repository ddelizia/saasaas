package db

import "github.com/ddelizia/saasaas/pkg/errors"

const (
	DataCreationError errors.TraceableType = "DataCreationError"
	DataUpdateError   errors.TraceableType = "DataUpdateError"
	DataGetError      errors.TraceableType = "DataGetError"
	DataListError     errors.TraceableType = "DataListError"
	DataDeleteError   errors.TraceableType = "DataDeleteError"
)

const (
	ModelFieldAccountID string = "AccountID"
	ModelFieldUserID    string = "UserID"
	ModelFieldAppID     string = "AppID"
)
