package repository

import (
	"billing-engine/domain/models"
	"billing-engine/infra"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type BillingConfig struct {
	db  *infra.DatabaseList
	log *logrus.Logger
}

func newBilling(db *infra.DatabaseList, logger *logrus.Logger) BillingConfig {
	return BillingConfig{
		db:  db,
		log: logger,
	}
}

type Billing interface {
	GetListUsers(ctx context.Context) ([]models.UserResponse, error)
	CheckExistUsers(ctx context.Context, userId int) (bool, error)
	CreateUserLoan(ctx context.Context, tx *sql.Tx, body models.LoanRequest) (int, error)
	CreatePaymentSchedule(ctx context.Context, tx *sql.Tx, body models.PaymentLoanRequest) error
	CheckIsDelinquents(ctx context.Context, param models.BillingRequest) (*models.DelinquentResponse, error)
	GetUserBillings(ctx context.Context, param models.BillingRequest) (*models.BillingResponse, error)
	CheckExistLoans(ctx context.Context, userId, loanId int) (bool, error)
	GetListLoans(ctx context.Context, body models.PaymentRequest) ([]models.PaymentResponse, error)
	UpdateUserWeeklyLoan(ctx context.Context, tx *sql.Tx, body models.PaymentRequest) error
	CountPendingLoans(ctx context.Context, tx *sql.Tx, body models.PaymentRequest) (int, error)
	UpdateUserLoan(ctx context.Context, tx *sql.Tx, body models.PaymentRequest) error
	GetUserSchedules(ctx context.Context, param models.BillingRequest) (*models.ScheduleResponse, error)
}

func (bc BillingConfig) GetListUsers(ctx context.Context) ([]models.UserResponse, error) {
	var result []models.UserResponse

	q := `SELECT id, name FROM users`

	err := bc.db.Backend.Read.SelectContext(ctx, &result, q)
	if err != nil && err != sql.ErrNoRows {
		return result, err
	}

	return result, nil
}

func (bc BillingConfig) CheckExistUsers(ctx context.Context, userId int) (bool, error) {
	var result bool

	q := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	err := bc.db.Backend.Read.GetContext(ctx, &result, q, userId)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (bc BillingConfig) CreateUserLoan(ctx context.Context, tx *sql.Tx, body models.LoanRequest) (int, error) {
	var result int

	q := `INSERT INTO loans
    (user_id, amount, interest, week, total, status, weekly_payment, start_date, created_at, updated_at)
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;
    `

	now := time.Now().UTC()

	err := tx.QueryRowContext(
		ctx,
		q,
		body.UserID,
		body.Amount,
		body.InterestRate,
		body.TotalWeeks,
		body.TotalAmount,
		body.Status,
		body.WeeklyPayment,
		body.StartDate,
		now,
		now,
	).Scan(&result)

	if err != nil {
		return 0, err
	}

	return result, nil
}

func (bc BillingConfig) CreatePaymentSchedule(ctx context.Context, tx *sql.Tx, body models.PaymentLoanRequest) error {
	for week := 1; week <= body.TotalWeeks; week++ {
		dueDate := body.StartDate.AddDate(0, 0, (week-1)*7) // Every 7 days
		now := time.Now().UTC()

		q := `INSERT INTO payments 
              (loan_id, user_id, week_number, amount, due_date, status, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

		_, err := tx.ExecContext(ctx, q, body.LoanID, body.UserID, week, body.Amount, dueDate, body.Status, now, now)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bc BillingConfig) CheckIsDelinquents(ctx context.Context, param models.BillingRequest) (*models.DelinquentResponse, error) {
	var result models.DelinquentResponse

	q := `SELECT week_number
		FROM payments
		WHERE user_id = $1 and loan_id = $2 AND status = 1 AND due_date < $3
		ORDER BY week_number`

	rows, err := bc.db.Backend.Read.QueryContext(ctx, q, param.UserID, param.LoanID, time.Now().UTC())
	if err != nil {
		return &result, err
	}
	defer rows.Close()

	var missedPayments, lastWeek int

	for rows.Next() {
		var weekNumber int
		if err := rows.Scan(&weekNumber); err != nil {
			return &result, err
		}

		if lastWeek != 0 && weekNumber == lastWeek+1 {
			missedPayments++
		} else {
			missedPayments = 1
		}
		lastWeek = weekNumber

		if missedPayments >= 2 {
			result.IsDelinquent = true
			return &result, nil
		}
	}

	if err := rows.Err(); err != nil {
		return &result, err
	}

	return &result, nil
}

func (bc BillingConfig) GetUserBillings(ctx context.Context, param models.BillingRequest) (*models.BillingResponse, error) {
	var result models.BillingResponse

	q := `SELECT total as amount_total FROM loans 
	WHERE user_id = $1 and id = $2`

	err := bc.db.Backend.Read.GetContext(ctx, &result.AmountTotal, q, param.UserID, param.LoanID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	q2 := `SELECT COALESCE(SUM(amount), 0) as amount_left FROM payments
	WHERE user_id = $1 and loan_id = $2 AND status = 1`

	err = bc.db.Backend.Read.GetContext(ctx, &result.AmountLeft, q2, param.UserID, param.LoanID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &result, nil
}

func (bc BillingConfig) CheckExistLoans(ctx context.Context, userId, loanId int) (bool, error) {
	var result bool

	q := `SELECT EXISTS(SELECT 1 FROM loans WHERE user_id = $1 and id = $2)`

	err := bc.db.Backend.Read.GetContext(ctx, &result, q, userId, loanId)
	if err != nil {
		return false, err
	}

	return result, nil
}

func (bc BillingConfig) GetListLoans(ctx context.Context, body models.PaymentRequest) ([]models.PaymentResponse, error) {
	var result []models.PaymentResponse

	q := `SELECT id, amount FROM payments where user_id = $1 and loan_id = $2 and status = 1 ORDER BY id ASC LIMIT $3`

	err := bc.db.Backend.Read.SelectContext(ctx, &result, q, body.UserID, body.LoanID, body.Week)
	if err != nil && err != sql.ErrNoRows {
		return result, err
	}

	if err == sql.ErrNoRows {
		return result, errors.New("no pending payments found for this loan")
	}

	return result, nil
}

func (bc BillingConfig) UpdateUserWeeklyLoan(ctx context.Context, tx *sql.Tx, body models.PaymentRequest) error {
	now := time.Now().UTC()

	q := `UPDATE payments SET status = 2, paid_date = $1, updated_at = $2 WHERE id = ANY($3)`

	_, err := tx.ExecContext(ctx, q, now, now, pq.Array(body.ListLoanID))
	if err != nil {
		return err
	}

	return nil
}

func (bc BillingConfig) CountPendingLoans(ctx context.Context, tx *sql.Tx, body models.PaymentRequest) (int, error) {
	var result int

	q := `SELECT count(id) FROM payments where user_id = $1 and loan_id = $2 and status = 1`

	err := tx.QueryRowContext(ctx, q, body.UserID, body.LoanID).Scan(&result)
	if err != nil && err != sql.ErrNoRows {
		return result, err
	}

	return result, nil
}

func (bc BillingConfig) UpdateUserLoan(ctx context.Context, tx *sql.Tx, body models.PaymentRequest) error {
	now := time.Now().UTC()

	q := `UPDATE loans SET status = 2, updated_at = $1 WHERE user_id = $2 and id = $3`

	_, err := tx.ExecContext(ctx, q, now, body.UserID, body.LoanID)
	if err != nil {
		return err
	}

	return nil
}

func (bc BillingConfig) GetUserSchedules(ctx context.Context, param models.BillingRequest) (*models.ScheduleResponse, error) {
	var result models.ScheduleResponse

	q := `SELECT amount, week, interest, weekly_payment ,start_date, status FROM loans where user_id = $1 and id = $2`

	err := bc.db.Backend.Read.GetContext(ctx, &result, q, param.UserID, param.LoanID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	q2 := `SELECT amount, week_number, status, due_date, COALESCE(TO_CHAR(paid_date, 'YYYY-MM-DD HH24:MI:SS'), '-') AS paid_date FROM payments where user_id = $1 and loan_id = $2 ORDER BY id ASC`

	err = bc.db.Backend.Read.SelectContext(ctx, &result.WeeklySchedule, q2, param.UserID, param.LoanID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &result, nil
}
