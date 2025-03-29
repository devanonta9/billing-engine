package repository

import (
	"billing-engine/domain/models"
	"billing-engine/infra"

	"github.com/sirupsen/logrus"
)

type Database struct {
	Billing Billing
}

type Repo struct {
	Database Database
}

func NewDatabase(db *infra.DatabaseList, logger *logrus.Logger) Database {
	return Database{
		Billing: newBilling(db, logger),
	}
}

func NewRepo(db *infra.DatabaseList, conf models.AppService, logger *logrus.Logger) Repo {
	return Repo{
		Database: NewDatabase(db, logger),
	}
}
