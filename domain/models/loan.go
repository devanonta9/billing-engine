package models

import "time"

type LoanRequest struct {
	UserID        int       `json:"user_id"`
	Amount        float64   `json:"amount"`
	InterestRate  float64   `json:"interest"`
	TotalWeeks    int       `json:"week"`
	StartDate     time.Time `json:"start"`
	Status        int
	WeeklyPayment float64
	TotalAmount   float64
}

type PaymentLoanRequest struct {
	LoanID     int
	UserID     int
	Amount     float64
	TotalWeeks int
	StartDate  time.Time
	Status     int
}

type BillingRequest struct {
	UserID int
	LoanID int
}

type BillingResponse struct {
	AmountTotal float64 `json:"amount_total" db:"amount_total"`
	AmountLeft  float64 `json:"amount_left" db:"amount_left"`
}

type DelinquentResponse struct {
	IsDelinquent bool `json:"is_delinquent" db:"is_delinquent"`
}

type PaymentRequest struct {
	UserID     int     `json:"user_id"`
	LoanID     int     `json:"loan_id"`
	Amount     float64 `json:"amount"`
	Week       int     `json:"week"`
	ListLoanID []int
	PaidAt     time.Time
}

type PaymentResponse struct {
	ID     int     `db:"id"`
	Amount float64 `db:"amount"`
}

type ScheduleResponse struct {
	Amount         float64                  `json:"amount" db:"amount"`
	TotalWeek      int                      `json:"week" db:"week"`
	Interest       float64                  `json:"interest_percentage" db:"interest"`
	TotalAmount    float64                  `json:"total_amount"`
	WeeklyPayment  float64                  `json:"weekly_payment" db:"weekly_payment"`
	StartLoan      string                   `json:"start_date" db:"start_date"`
	Status         int                      `json:"-" db:"status"`
	StatusStr      string                   `json:"status"`
	WeeklySchedule []WeeklyScheduleResponse `json:"weekly_data"`
}

type WeeklyScheduleResponse struct {
	WeeklyAmount float64 `json:"amount_weekly" db:"amount"`
	WeekNumber   int     `json:"week_number" db:"week_number"`
	DueData      string  `json:"due_date" db:"due_date"`
	PaidDate     string  `json:"paid_date" db:"paid_date"`
	Status       int     `json:"-" db:"status"`
	StatusStr    string  `json:"status_weekly"`
}
