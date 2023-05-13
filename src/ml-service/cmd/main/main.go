package main

import (
	"ml/internal/config"
	"ml/internal/repository/postgres"
	"ml/internal/service"
	"ml/pkg/logging"
	"ml/pkg/psqlcli"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()

	cfg := config.GetConfig()
	cli, err := psqlcli.New(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode)
	if err != nil {
		logger.Fatal(err)
	}

	//router := gin.New()
	//router.LoadHTMLGlob("etc/static/*.html")
	//
	sampleRepo := postgres.NewSamplePostgres(cli)
	frameRepo := postgres.NewFramePostgres(cli)
	//
	sampleService := service.NewSampleService(sampleRepo, frameRepo)
	mlService := service.NewMLService(logger, sampleService)

	actual, predicted, err := mlService.Test()
	if err != nil {
		panic(err)
	}

	mlService.GetHeatmap(actual, predicted)
	//
	//mlHandler := handler.NewMLHandler(logger, sampleService, mlService)
	//mlHandler.Register(router)
	//
	//logger.Info("ml handler registered")
	//
	//server.Run(cfg, router, logger)
}
