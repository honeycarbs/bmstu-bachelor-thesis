package main

import (
	"cluster/cmd/cli"
	"fmt"
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
	//cli, err := psqlcli.NewClient(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode)
	//if err != nil {
	//	logger.Fatal(err)
	//}
	//
	//frameRepo := postgres.NewFramePostgres(cli)
	//clusterRepo := postgres.NewClusterPostgres(cli)
	//
	//frameService := service.NewFrameService(frameRepo)
	//
	//frames, err := frameService.GetAll()
	//if err != nil {
	//	panic(err)
	//}
	//
	//clusterService := service.NewClusterService(clusterRepo, logger)
	//
	//logger.Infof("Collected %v frames", len(frames))
	//clusterAmount := config.GetConfig().ClusterAmount
	//maxRounds := config.GetConfig().KMeansMaxRounds
	//
	//logger.Infof("creating %v kmeans clusters with %v iterations", clusterAmount, maxRounds)
	//clusters, err := clusterService.AssignClusters(frames, clusterAmount, maxRounds)
	//if err != nil {
	//	logger.Info(err)
	//	panic(err)
	//}
	//
	//logger.Info("Started creating cluster entities...")
	//for i, cluster := range clusters {
	//	err := clusterService.CreateCluster(cluster)
	//	logger.Infof("Created cluster entity no. (%v)", i)
	//	if err != nil {
	//		logger.Info(err)
	//		panic(err)
	//	}
	//}
	//
	//logger.Info("Started assigning clusters to frames...")
	//for i, fm := range frames {
	//	err := frameService.AssignCluster(fm, clusters)
	//	if err != nil {
	//		logger.Info(err)
	//		panic(err)
	//	}
	//	logger.Infof("Assigned cluster to frame (%v) out of (%v)", i, len(frames))
	//}
	//logger.Info("All the clusters are assigned, finishing...")
}
