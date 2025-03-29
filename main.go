package main

import (
	"fmt"
	"net/http"

	"billing-engine/domain/models"
	"billing-engine/handler"
	"billing-engine/infra"
	"billing-engine/repository"
	"billing-engine/routes"
	"billing-engine/service"

	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	conf, err := getConfigKey()
	if err != nil {
		panic(err)
	}

	handler, log, err := newRepoContext(conf)
	if err != nil {
		panic(err)
	}

	headers := handlers.AllowedHeaders(conf.Route.Headers)
	methods := handlers.AllowedMethods(conf.Route.Methods)
	origins := handlers.AllowedOrigins(conf.Route.Origins)
	credentials := handlers.AllowCredentials()

	router := routes.GetCoreEndpoint(conf, handler, log)

	port := fmt.Sprintf(":%s", conf.App.Port)
	log.Info("server listen to port ", port)
	log.Fatal(http.ListenAndServe(port, handlers.CORS(headers, methods, origins, credentials)(router)))
}

func getConfigKey() (*models.AppService, error) {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	conf := models.AppService{
		App: models.App{
			Name: viper.GetString("APP.NAME"),
			Port: viper.GetString("APP.PORT"),
		},
		Route: models.Route{
			Methods: viper.GetStringSlice("ROUTE.METHODS"),
			Headers: viper.GetStringSlice("ROUTE.HEADERS"),
			Origins: viper.GetStringSlice("ROUTE.ORIGIN"),
		},
		Database: models.Database{
			Read: models.DBDetail{
				Username:     viper.GetString("DATABASE.READ.USERNAME"),
				Password:     viper.GetString("DATABASE.READ.PASSWORD"),
				URL:          viper.GetString("DATABASE.READ.URL"),
				Port:         viper.GetString("DATABASE.READ.PORT"),
				DBName:       viper.GetString("DATABASE.READ.DB_NAME"),
				MaxIdleConns: viper.GetInt("DATABASE.READ.MAXIDLECONNS"),
				MaxOpenConns: viper.GetInt("DATABASE.READ.MAXOPENCONNS"),
				MaxLifeTime:  viper.GetInt("DATABASE.READ.MAXLIFETIME"),
				Timeout:      viper.GetString("DATABASE.READ.TIMEOUT"),
				SSLMode:      viper.GetString("DATABASE.READ.SSL_MODE"),
			},
		},
	}

	return &conf, nil
}

func newRepoContext(conf *models.AppService) (handler.Handler, *logrus.Logger, error) {
	var handlers handler.Handler

	// Init log
	logger := infra.NewLogger(conf)

	// Init DB Write Connection.
	dbRead := infra.NewDB(logger)
	dbRead.ConnectDB(&conf.Database.Read)
	if dbRead.Err != nil {
		return handlers, logger, dbRead.Err
	}

	dbList := &infra.DatabaseList{
		Backend: infra.DatabaseType{
			Read: &dbRead,
		},
	}

	// Init Minio config.
	repo := repository.NewRepo(dbList, *conf, logger)
	usecase := service.NewService(repo, *conf, dbList, logger)
	handlers = handler.NewHandler(usecase, *conf, logger)

	return handlers, logger, nil
}
