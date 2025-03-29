package routes

import (
	"billing-engine/domain/models"
	"billing-engine/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func GetCoreEndpoint(conf *models.AppService, handler handler.Handler, log *logrus.Logger) *mux.Router {
	initRoute := mux.NewRouter()
	nonJWTRoute := initRoute.PathPrefix("").Subrouter()

	// Get Endpoint.
	getV1(nonJWTRoute, conf, handler)

	return initRoute
}

func getV1(router *mux.Router, conf *models.AppService, handler handler.Handler) {
	router.HandleFunc("/v1/users", handler.Billing.GetUsers).Methods(http.MethodGet)
	router.HandleFunc("/v1/loans", handler.Billing.CreateLoans).Methods(http.MethodPost)
	router.HandleFunc("/v1/delinquents", handler.Billing.CheckDelinquents).Methods(http.MethodGet)
	router.HandleFunc("/v1/billings", handler.Billing.GetBillings).Methods(http.MethodGet)
	router.HandleFunc("/v1/payments", handler.Billing.CreatePayments).Methods(http.MethodPost)
	router.HandleFunc("/v1/schedules", handler.Billing.GetSchedules).Methods(http.MethodGet)
}
