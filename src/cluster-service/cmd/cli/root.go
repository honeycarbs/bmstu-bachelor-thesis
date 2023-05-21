package cli

import (
	"cluster/internal/config"
	"cluster/internal/repository/postgres"
	"cluster/internal/service"
	"cluster/pkg/logging"
	"cluster/pkg/psqlcli"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	mode string
	path string

	RootCmd = &cobra.Command{
		Use:   "cluster-service",
		Short: "Provides k-means clustering function interface",
		Long: `Part of BMSTU batchelor thesis for SER
			Using HMM and MFCC.
			Configuration file is provided in ./etc`,
		RunE: commandHandler,
	}
)

func init() {
	RootCmd.PersistentFlags().StringVar(&mode, "mode", "", "cluster all or add to existing")
	RootCmd.PersistentFlags().StringVar(&path, "path", "", "absolute path for file to cluster")
}

func commandHandler(cmd *cobra.Command, args []string) error {
	logging.Init()
	logger := logging.GetLogger()

	cfg := config.GetConfig()
	cli, err := psqlcli.NewClient(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode)
	if err != nil {
		logger.Fatal(err)
	}

	frameRepo := postgres.NewFramePostgres(cli)
	frameService := service.NewFrameService(frameRepo)

	clusterRepo := postgres.NewClusterPostgres(cli)
	clusterService := service.NewClusterService(clusterRepo, logger)

	logger.Infof("mode = %v", mode)

	switch mode {
	case "assign":
		{
			logger.Info("Got assign command, executing")

			frames, err := frameService.GetAll()
			if err != nil {
				panic(err)
			}

			logger.Infof("Collected %v frames", len(frames))
			clusterAmount := config.GetConfig().ClusterAmount
			maxRounds := config.GetConfig().KMeansMaxRounds

			logger.Infof("creating %v kmeans clusters with %v iterations", clusterAmount, maxRounds)
			clusters, err := clusterService.AssignClusters(frames, clusterAmount, maxRounds)
			if err != nil {
				logger.Info(err)
				panic(err)
			}

			logger.Info("Started creating cluster entities...")
			for i, cluster := range clusters {
				err := clusterService.CreateCluster(cluster)
				logger.Infof("Created cluster entity no. (%v)", i)
				if err != nil {
					logger.Info(err)
					panic(err)
				}
			}

			err = frameService.AssignClusters(frames, clusters)
			if err != nil {
				logger.Info(err)
				panic(err)
			}

			logger.Info("All the clusters are assigned, finishing...")
		}
	case "add":
		{
			logger.Info("Got add command, executing")
			frames, err := frameService.GetOne(path)
			if err != nil {
				return err
			}

			logger.Infof("Got %v frames", len(frames))

			clusters, err := clusterService.GetAll()

			err = frameService.AssignClusters(frames, clusters)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown command")
	}
	return nil
}
