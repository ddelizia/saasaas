package main

import (
	"net/http"

	"github.com/ddelizia/saasaas/pkg/account"
	"github.com/ddelizia/saasaas/pkg/account/dynamodb"
	"github.com/ddelizia/saasaas/pkg/db/dynamodbsaas"
	"github.com/ddelizia/saasaas/pkg/transport/rest"
	transporthttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func main() {

	tableName := dynamodbsaas.TableNameWithEnv("account")
	db := dynamodbsaas.DynamoDbInstance()

	repo := dynamodb.NewRepository(db, tableName)
	svc := account.New(repo)

	routerFunc := func(r *mux.Router) {

		var options []transporthttp.ServerOption

		r.Methods(http.MethodPost).Path("/account/").Handler(transporthttp.NewServer(
			account.CreateEP(svc),
			rest.HttpRequestDecoder(&account.Account{}),
			rest.HttpResponseEncoder(),
			options...,
		))

		r.Methods(http.MethodGet).Path("/account/{id}").Handler(transporthttp.NewServer(
			account.GetEP(svc),
			rest.HttpRequestDecoder(nil),
			rest.HttpResponseEncoder(),
			options...,
		))

		r.Methods(http.MethodDelete).Path("/account/{id}").Handler(transporthttp.NewServer(
			account.DeleteEP(svc),
			rest.HttpRequestDecoder(nil),
			rest.HttpResponseEncoder(),
			options...,
		))

		r.Methods(http.MethodPut).Path("/account/").Handler(transporthttp.NewServer(
			account.UpdateEP(svc),
			rest.HttpRequestDecoder(&account.Account{}),
			rest.HttpResponseEncoder(),
			options...,
		))
	}

	svr := rest.HttpServer(routerFunc, 8080)
	svr.ListenAndServe()

}
