package constants

import "time"

const (
	HandlerErrorAuthInvalid                string = "authorization invalid"
	HandlerErrorAuthInvalidID              string = "authorization tidak valid"
	HandlerErrorResponseKeyIDEmpty         string = "key id cannot be empty"
	HandlerErrorRequestDataNotValid        string = "request data not valid"
	HandlerErrorRequestDataNotValidID      string = "data request tidak valid"
	HandlerErrorRequestDataEmpty           string = "request data empty"
	HandlerErrorRequestDataEmptyID         string = "data request kosong"
	HandlerErrorRequestDataFormatInvalid   string = "request data format invalid"
	HandlerErrorRequestDataFormatInvalidID string = "format data request salah"
	HandlerErrorCookiesEmpty               string = "key data cannot be empty"
	HandlerErrorCookiesInvalid             string = "key data invalid"
	HandlerErrorKeyIDInvalid               string = "key id invalid"
	HandlerErrorImageSizeTooLarge          string = "image too large, max size 1 Mb"
	HandlerErrorImageSizeTooLargeID        string = "file terlalu besar, maks. 1 Mb"
	HandlerErrorImageDataInvalid           string = "image data invalid"
	HandlerErrorImageDataInvalidID         string = "data gambar tidak sesuai"
	HandlerErrorImageDataEmpty             string = "image data cannot be empty"
	HandlerErrorImageDataEmptyID           string = "data gambar tidak boleh kosong"
	HandlerErrorFileSizeTooLarge           string = "file too large, max size 1 Mb"
	HandlerErrorFileDataInvalid            string = "file data invalid"
	HandlerErrorFileDataEmpty              string = "file data cannot be empty"
)

const (
	ConnectDBSuccess string = "Connected to DB"

	ConnectDBFail string = "Could not connect database, error"

	ClosingDBSuccess string = "Database conn gracefully close"
	ClosingDBFailed  string = "Error closing DB connection"

	Success string = "success"
	Fail    string = "fail"

	DataNotFound string = "no data found"

	DBTimeLayout       string = "2006-01-02 15:04:05"
	ResponseTimeLayout string = "2006-01-02T15:04:05-0700"
)

const (
	FullTimeFormat        string = "2006-01-02 15:04:05"
	DisplayDateTimeFormat string = "02 Jan 2006 15:04:05"
	DateFormat            string = "2006-01-02"
)

const (
	LogRotationTime = time.Duration(24) * time.Hour
	MaxRotationFile = 4
)

const (
	LoanStatusOpen          = 1
	LoanStatusClose         = 2
	LoanWeeklyStatusPending = 1
	LoanWeeklYStatusPaid    = 2
)
