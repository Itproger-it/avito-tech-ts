package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	stng "micro-service"
	"micro-service/pkg/handler"
	rp "micro-service/pkg/repository"
	"micro-service/pkg/service"
	database "micro-service/storage"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	viper.AutomaticEnv()
	// if err := initConfig(); err != nil {
	// 	logrus.Fatalf("error run config: %s", err.Error())
	// }

	db, err := database.NewStorage(*database.NewConfig())
	defer func() {
		db.Close()
	}()
	if err != nil {
		logrus.Fatalf("error run database: %s", err.Error())
	}

	if err != nil {
		logrus.Fatalf("eroor create snowflake node: %s", err.Error())
	}

	rp := rp.NewRepository(db)
	service := service.NewService(rp)
	handler := handler.NewHandler(service)

	logrus.Infof("Starting server on port %s", "http://localhost:")
	server := stng.Server{}
	if err := server.Run(viper.GetString("BACKEND_PORT"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("error run server: %s", err.Error())
	}
}

// func initConfig() error {
// 	viper.SetConfigFile("../config/config.yaml")
// 	return viper.ReadInConfig()
// }
