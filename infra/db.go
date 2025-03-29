package infra

import (
	"context"
	"database/sql"
	"os"
	"time"

	"billing-engine/domain/constants"
	"billing-engine/domain/models"
	"billing-engine/domain/utils"

	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	_ "github.com/lib/pq"
	"github.com/rifflock/lfshook"
)

// IDatabase is interface for database
type Database interface {
	ConnectDB(dbAcc *models.DBDetail)
	Close()

	Begin() (*sql.Tx, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type DatabaseList struct {
	Backend DatabaseType
}

type DatabaseType struct {
	Read Database
}

// DBHandler - Database struct.
type DBHandler struct {
	DB  *sqlx.DB
	Err error
	log *log.Logger
}

func NewDB(log *log.Logger) DBHandler {
	return DBHandler{
		log: log,
	}
}

// ConnectDB - function for connect DB.
func (d *DBHandler) ConnectDB(db *models.DBDetail) {
	dbs, err := sqlx.Open("postgres", "user="+db.Username+" password="+db.Password+" sslmode="+db.SSLMode+" dbname="+db.DBName+" host="+db.URL+" port="+db.Port+" connect_timeout="+db.Timeout)
	if err != nil {
		log.Error(constants.ConnectDBFail, err.Error())
		d.Err = err
	}

	d.DB = dbs

	err = d.DB.Ping()
	if err != nil {
		log.Error(constants.ConnectDBFail, err.Error())
		d.Err = err
	}

	d.log.Info(constants.ConnectDBSuccess)
	d.DB.SetConnMaxLifetime(time.Duration(db.MaxLifeTime))
}

// Close - function for connection lost.
func (d *DBHandler) Close() {
	if err := d.DB.Close(); err != nil {
		d.log.Println(constants.ClosingDBFailed, err.Error())
	} else {
		d.log.Println(constants.ClosingDBSuccess)
	}
}

func (d *DBHandler) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	result, err := d.DB.QueryContext(ctx, query, args...)
	return result, err
}

func (d *DBHandler) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := d.DB.SelectContext(ctx, dest, query, args...)
	return err
}

func (d *DBHandler) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := d.DB.GetContext(ctx, dest, query, args...)
	return err
}

func (d *DBHandler) Begin() (*sql.Tx, error) {
	return d.DB.Begin()
}

var logger *log.Logger

func NewLogger(conf *models.AppService) *log.Logger {
	if logger == nil {
		path := "log/"

		isExist, err := utils.DirExists(path)
		if err != nil {
			panic(err)
		}

		if !isExist {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		writer, err := rotatelogs.New(
			path+conf.App.Name+"-"+"%Y%m%d.log",
			rotatelogs.WithMaxAge(-1),
			rotatelogs.WithRotationCount(constants.MaxRotationFile),
			rotatelogs.WithRotationTime(constants.LogRotationTime),
		)
		if err != nil {
			panic(err)
		}

		logger = log.New()

		// Set Hook with writer & formatter for log file
		logger.Hooks.Add(lfshook.NewHook(
			writer,
			&log.TextFormatter{
				DisableColors:   false,
				FullTimestamp:   true,
				TimestampFormat: constants.FullTimeFormat,
			},
		))

		// Set formatter for os.Stdout
		logger.SetFormatter(&log.TextFormatter{
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: constants.FullTimeFormat,
		})

		return logger
	}

	return logger
}
