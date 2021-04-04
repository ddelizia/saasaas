package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/ddelizia/saasaas/pkg/errors"
	"github.com/ddelizia/saasaas/pkg/t"
	transporthttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func bearerToken(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) < 2 {
		return ""
	}
	return splitToken[1]
}

func apiKey(r *http.Request) string {
	return r.Header.Get("X-API-Key")
}

func HttpRequestDecoder(input interface{}) transporthttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var in interface{} = nil
		if input != nil {
			in := t.GenerareEmptyInterface(input)
			if err := json.NewDecoder(r.Body).Decode(in); err != nil {
				return nil, err
			}
		}

		vars := mux.Vars(r)
		req := &GenericRequest{
			Method:      r.Method,
			Body:        in,
			Path:        r.URL.Path,
			QueryParams: r.URL.Query(),
			PathParams:  vars,
			BearerToken: bearerToken(r),
			ApiKey:      apiKey(r),
		}
		logrus.Info(fmt.Sprintf("Request [%s] %s - ? %v, p %v", req.Method, req.Path, req.QueryParams, req.PathParams))
		return req, nil
	}
}

func HttpResponseEncoder() transporthttp.EncodeResponseFunc {
	return func(_ context.Context, w http.ResponseWriter, response interface{}) error {
		r := response.(*GenericResponse)
		w.Header().Add("Content-Type", "application/json")

		if r.Err != nil {
			var status *ErrorData
			if reflect.TypeOf(r.Err) == reflect.TypeOf(errors.Traceable{}) {
				e := r.Err.(errors.Traceable)
				status = r.ErrorStatusMapping[e.ErrorType]
				if status == nil {
					status = &ErrorData{http.StatusInternalServerError, "Server error: not mapped"}
				}
			} else {
				status = &ErrorData{http.StatusInternalServerError, "Server error: unknown"}
			}
			w.WriteHeader(status.Status)
			err := json.NewEncoder(w).Encode(status)
			return err
		}

		if r.Body != nil {
			err := json.NewEncoder(w).Encode(r.Body)
			if err != nil {
				status := &ErrorData{http.StatusInternalServerError, "Server error: encoding"}
				w.WriteHeader(status.Status)
				err := json.NewEncoder(w).Encode(status)
				return err
			}
		}
		w.WriteHeader(r.Status)

		return nil
	}
}

type RouterFunc = func(*mux.Router)

func HttpServer(routerFunc RouterFunc, port int) *http.Server {
	r := mux.NewRouter()

	routerFunc(r)

	logrus.Info("status listening port 8080")

	return &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(port),
		Handler: r,
	}
}
