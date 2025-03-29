package service

import (
	"billing-engine/domain/constants"
	"billing-engine/domain/models"
	"billing-engine/infra"
	"billing-engine/repository"
	"context"
	"errors"
	"log"

	"github.com/sirupsen/logrus"
)

type BillingService struct {
	db     repository.Database
	dbList *infra.DatabaseList
	log    *logrus.Logger
}

func newBillingService(database repository.Database, dbList *infra.DatabaseList, logger *logrus.Logger) BillingService {
	return BillingService{
		db:     database,
		dbList: dbList,
		log:    logger,
	}
}

type Billing interface {
	GetListUsers(ctx context.Context) ([]models.UserResponse, string, error)
	CreateUserLoan(ctx context.Context, body models.LoanRequest) (string, error)
	CheckIsDelinquents(ctx context.Context, param models.BillingRequest) (*models.DelinquentResponse, string, error)
	GetUserBillings(ctx context.Context, param models.BillingRequest) (*models.BillingResponse, string, error)
	CreateUserPayment(ctx context.Context, body models.PaymentRequest) (string, error)
	GetUserSchedules(ctx context.Context, param models.BillingRequest) (*models.ScheduleResponse, string, error)
}

func (bs BillingService) GetListUsers(ctx context.Context) ([]models.UserResponse, string, error) {
	users, err := bs.db.Billing.GetListUsers(ctx)
	if err != nil {
		log.Println("GetListUsers | GetListUsers | Medium", err.Error())
		return users, "failed to get list users", err
	}

	return users, "success", nil
}

func (bs BillingService) CreateUserLoan(ctx context.Context, body models.LoanRequest) (string, error) {
	exist, err := bs.db.Billing.CheckExistUsers(ctx, body.UserID)
	if err != nil {
		log.Println("CreateUserLoan | CheckExistUsers | Medium", err.Error())
		return "failed to check exist users", err
	}

	if !exist {
		log.Println("CreateUserLoan | user not exist | Medium")
		return "failed to check exist users", errors.New("user not exist")
	}

	tx, err := bs.dbList.Backend.Read.Begin()
	if err != nil {
		log.Println("CreateUserLoan | Begin TX | Low", err.Error())
		return "failed to begin transaction", err
	}
	defer tx.Rollback()

	interestAmount := body.Amount * body.InterestRate / 100
	body.TotalAmount = body.Amount + interestAmount
	body.WeeklyPayment = body.TotalAmount / float64(body.TotalWeeks)
	body.Status = 1

	loanId, err := bs.db.Billing.CreateUserLoan(ctx, tx, body)
	if err != nil {
		log.Println("CreateUserLoan | CreateUserLoan | High", err.Error())
		return "failed to create loan", err
	}

	if loanId > 0 {
		var bodyPayment models.PaymentLoanRequest
		bodyPayment.UserID = body.UserID
		bodyPayment.LoanID = loanId
		bodyPayment.StartDate = body.StartDate
		bodyPayment.TotalWeeks = body.TotalWeeks
		bodyPayment.Amount = body.WeeklyPayment
		bodyPayment.Status = 1

		err = bs.db.Billing.CreatePaymentSchedule(ctx, tx, bodyPayment)
		if err != nil {
			log.Println("CreateUserLoan | CreatePaymentSchedule | High", err.Error())
			return "failed to create payment schedule", err
		}
	}
	if err = tx.Commit(); err != nil {
		log.Println("CreateUserLoan | Commit TX | Low", err.Error())
		return "failed to commit transaction", err
	}

	return "success", nil
}

func (bs BillingService) CheckIsDelinquents(ctx context.Context, param models.BillingRequest) (*models.DelinquentResponse, string, error) {
	delinquent, err := bs.db.Billing.CheckIsDelinquents(ctx, param)
	if err != nil {
		log.Println("CheckIsDelinquents | CheckIsDelinquents | High", err.Error())
		return nil, "failed to get user delinquents", err
	}

	return delinquent, "success", nil
}

func (bs BillingService) GetUserBillings(ctx context.Context, param models.BillingRequest) (*models.BillingResponse, string, error) {
	billings, err := bs.db.Billing.GetUserBillings(ctx, param)
	if err != nil {
		log.Println("GetUserBillings | GetUserBillings | High", err.Error())
		return nil, "failed to get user billings", err
	}

	return billings, "success", nil
}

func (bs BillingService) CreateUserPayment(ctx context.Context, body models.PaymentRequest) (string, error) {
	var listLoan []int

	existUser, err := bs.db.Billing.CheckExistUsers(ctx, body.UserID)
	if err != nil {
		log.Println("CreateUserPayment | CheckExistUsers | Medium", err.Error())
		return "failed to check exist users", err
	}

	if !existUser {
		log.Println("CreateUserPayment | user not exist | Medium")
		return "failed to check exist users", errors.New("user not exist")
	}

	existLoan, err := bs.db.Billing.CheckExistLoans(ctx, body.UserID, body.LoanID)
	if err != nil {
		log.Println("CreateUserPayment | CheckExistLoans | Medium", err.Error())
		return "failed to check exist loans", err
	}

	if !existLoan {
		log.Println("CreateUserPayment | loan not exist | Medium")
		return "failed to check exist loans", errors.New("loan not exist")
	}

	tx, err := bs.dbList.Backend.Read.Begin()
	if err != nil {
		log.Println("CreateUserPayment | Begin TX | Low", err.Error())
		return "failed to begin transaction", err
	}
	defer tx.Rollback()

	loanId, err := bs.db.Billing.GetListLoans(ctx, body)
	if err != nil {
		log.Println("CreateUserPayment | GetUserLoans | High", err.Error())
		return "failed to get loan", err
	}

	if len(loanId) > 0 {
		for _, v := range loanId {
			listLoan = append(listLoan, v.ID)
		}
		body.ListLoanID = listLoan
		if len(body.ListLoanID) != body.Week {
			log.Println("CreateUserPayment | len(listLoan) != body.Week | Medium")
			return "failed to pay due to week and loan remaining not match", errors.New("week and loan remaining not match")
		}

		if loanId[0].Amount != body.Amount {
			log.Println("CreateUserPayment | loanId[0].Amount != body.Amount | Medium")
			return "failed to pay due to paid amount and weekly loan amount not match", errors.New("amount and weekly loan amount not match")
		}

		err = bs.db.Billing.UpdateUserWeeklyLoan(ctx, tx, body)
		if err != nil {
			log.Println("CreateUserPayment | UpdateUserWeeklyLoan | High", err.Error())
			return "failed to update weekly loan", err
		}

		count, err := bs.db.Billing.CountPendingLoans(ctx, tx, body)
		if err != nil {
			log.Println("CreateUserPayment | CountPendingLoans | Medium", err.Error())
			return "failed to count loan", err
		}

		if count == 0 {
			err = bs.db.Billing.UpdateUserLoan(ctx, tx, body)
			if err != nil {
				log.Println("CreateUserPayment | UpdateUserLoan | High", err.Error())
				return "failed to update loan", err
			}
		}
	}
	if err = tx.Commit(); err != nil {
		log.Println("CreateUserPayment | Commit TX | Low", err.Error())
		return "failed to commit transaction", err
	}

	return "success", nil
}

func (bs BillingService) GetUserSchedules(ctx context.Context, param models.BillingRequest) (*models.ScheduleResponse, string, error) {
	schedules, err := bs.db.Billing.GetUserSchedules(ctx, param)
	if err != nil {
		log.Println("GetUserSchedules | GetUserSchedules | High", err.Error())
		return nil, "failed to get user schedules", err
	}

	var LoanStatusDescription = map[int]string{
		constants.LoanStatusOpen:  "Open",
		constants.LoanStatusClose: "Close",
	}

	var LoanWeeklyStatusDescription = map[int]string{
		constants.LoanWeeklyStatusPending: "Pending",
		constants.LoanWeeklYStatusPaid:    "Paid",
	}

	if schedules != nil {
		interestAmount := schedules.Amount * schedules.Interest / 100
		schedules.TotalAmount = schedules.Amount + interestAmount
		schedules.StatusStr = LoanStatusDescription[schedules.Status]
		for k, v := range schedules.WeeklySchedule {
			v.StatusStr = LoanWeeklyStatusDescription[v.Status]
			schedules.WeeklySchedule[k] = v
		}
	}

	return schedules, "success", nil
}
