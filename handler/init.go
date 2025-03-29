package handler

import (
	"billing-engine/domain/models"
	"billing-engine/service"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	Billing BillingHandler
}

func NewHandler(sv service.Service, conf models.AppService, logger *logrus.Logger) Handler {
	return Handler{
		Billing: newBillingHandler(sv.Billing, conf, logger),
	}
}
