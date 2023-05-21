package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"ml/internal/config"
	"ml/internal/pkg/heatmap"
	"ml/internal/repository/postgres"
	"ml/internal/service"
	"ml/pkg/logging"
	"ml/pkg/psqlcli"
)

var (
	mode string
	path string

	RootCmd = &cobra.Command{
		Use:   "ml-service",
		Short: "Provides HMM functions for learning",
		Long: `Part of BMSTU batchelor thesis for SER
			Using HMM and MFCC.
			Configuration file is provided in ./etc`,
		RunE: commandHandler,
	}
)

func init() {
	RootCmd.PersistentFlags().StringVar(&mode, "mode", "", "test, train or recognize mode for HMM")
	RootCmd.PersistentFlags().StringVar(&path, "path", "", "absolute path for file to recognize")
}

func commandHandler(cmd *cobra.Command, args []string) error {
	logging.Init()
	logger := logging.GetLogger()

	cfg := config.GetConfig()
	cli, err := psqlcli.New(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode)
	if err != nil {
		logger.Fatal(err)
	}

	sampleRepo := postgres.NewSamplePostgres(cli, logger)
	frameRepo := postgres.NewFramePostgres(cli)

	sampleService := service.NewSampleService(sampleRepo, frameRepo, logger)
	mlService := service.NewMLService(logger, sampleService)

	switch mode {
	case "test":
		{
			logger.Info("Got test command, executing")
			actual, predicted, err := mlService.Test()
			if err != nil {
				return err
			}

			heatmap.GetHeatmap(actual, predicted)
		}
	case "train":
		{
			err := mlService.Train()

			return err
		}
	case "recognize":
		{
			emotion, err := mlService.Recognize(path)
			if err != nil {
				return err
			}
			fmt.Println(emotion)

			return nil
		}
	default:
		return fmt.Errorf("unknown command")
	}
	return nil
}
