package account

import (
	"context"
	"net/http"

	"github.com/ddelizia/saasaas/pkg/errors"
	"github.com/ddelizia/saasaas/pkg/t"
	"github.com/ddelizia/saasaas/pkg/transport/rest"
	"github.com/go-kit/kit/endpoint"
	"github.com/sirupsen/logrus"
)

func CreateEP(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logrus.Info("Request type")
		req := request.(*rest.GenericRequest)
		body := req.Body.(*Account)
		resp, err := svc.Create(ctx, body)

		res := &rest.GenericResponse{
			Body: resp,
			Err:  err,
			ErrorStatusMapping: map[errors.TraceableType]*rest.ErrorData{
				CreationError: {http.StatusConflict, "Error while creating the account"},
			},
			Status: http.StatusCreated,
		}
		return res, nil
	}
}

func UpdateEP(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logrus.Info("Request type")
		req := request.(*rest.GenericRequest)
		body := req.Body.(*Account)
		id := req.PathParams["id"]
		err := svc.Update(ctx, t.NewString(id), body)

		res := &rest.GenericResponse{
			Err: err,
			ErrorStatusMapping: map[errors.TraceableType]*rest.ErrorData{
				UpdateError: {http.StatusConflict, "Error while updating the account"},
			},
			Status: http.StatusNoContent,
		}
		return res, nil
	}
}

func GetEP(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logrus.Info("Request type")
		req := request.(*rest.GenericRequest)
		id := req.PathParams["id"]
		acc, err := svc.Get(ctx, t.NewString(id))

		res := &rest.GenericResponse{
			Err: err,
			ErrorStatusMapping: map[errors.TraceableType]*rest.ErrorData{
				GetError: {http.StatusNotFound, "Not found"},
			},
			Body: acc,
		}
		return res, nil
	}
}

func DeleteEP(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		logrus.Info("Request type")
		req := request.(*rest.GenericRequest)
		id := req.PathParams["id"]
		err := svc.Delete(ctx, t.NewString(id))

		res := &rest.GenericResponse{
			Err: err,
			ErrorStatusMapping: map[errors.TraceableType]*rest.ErrorData{
				DeleteError: {http.StatusNotFound, "Not found"},
			},
			Status: http.StatusCreated,
		}
		return res, nil
	}
}
