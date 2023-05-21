package main

import (
	"fmt"
	"ml/cmd/cli"
	"os"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//logging.Init()
	//logger := logging.GetLogger()
	//
	//cfg := config.GetConfig()
	//cli, err := psqlcli.New(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode)
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//
	//sampleRepo := postgres.NewSamplePostgres(cli, logger)
	//frameRepo := postgres.NewFramePostgres(cli)
	//
	//sampleService := service.NewSampleService(sampleRepo, frameRepo, logger)
	//mlService := service.NewMLService(logger, sampleService)
	//
	//emo, err := mlService.Recognize(
	//	"/home/honeycarbs/BMSTU/bmstu-bachelor-thesis/src/DUSHA/crowd_train/wavs/d56b0a86113cb87896c870c04ea2f1db.wav")
	//
	//logger.Trace(emo)
}
