package rest

import (
	"github.com/ddelizia/saasaas/pkg/errors"
)

type ErrorData struct {
	Status  int
	Message string
}

type GenericResponse struct {
	Body               interface{}
	Status             int
	Err                error
	ErrorStatusMapping map[errors.TraceableType]*ErrorData
}

type GenericRequest struct {
	Method      string
	Path        string
	Body        interface{}
	QueryParams map[string][]string
	PathParams  map[string]string
	BearerToken string
	ApiKey      string
}
