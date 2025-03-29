package service

import (
	"billing-engine/domain/models"
	"billing-engine/infra"
	"billing-engine/repository"

	"github.com/sirupsen/logrus"
)

type Service struct {
	Billing BillingService
}

func NewService(repo repository.Repo, conf models.AppService, dbList *infra.DatabaseList, logger *logrus.Logger) Service {
	return Service{
		Billing: newBillingService(repo.Database, dbList, logger),
	}
}
