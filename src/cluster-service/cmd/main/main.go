package main

import (
	"cluster/internal/repository/postgres"
	"cluster/internal/service"
	"cluster/pkg/psqlcli"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.New()

func main() {
	logger.SetOutput(os.Stdout)
	cli, err := psqlcli.New("localhost", "5432", "admin", "admin", "thesis", "disable")
	if err != nil {
		panic(err)
	}

	sampleRepo := postgres.NewSamplePostgres(cli)
	frameRepo := postgres.NewFramePostgres(cli)
	clusterRepo := postgres.NewClusterPostgres(cli)

	sampleService := service.NewSampleService(sampleRepo)
	frameService := service.NewFrameService(frameRepo)
	clusterService := service.NewClusterService(clusterRepo)

	samples, err := sampleService.GetAll()
	if err != nil {
		panic(err)
	}

	//samples = samples[:100]

	for i := 0; i < len(samples); i++ {
		samples[i].Frames, err = frameService.GetAllBySample(samples[i].ID)
		if err != nil {
			panic(err)
		}
		logger.Infof("got %v frames from sample %v out of %v", len(samples[i].Frames), i, 100)
	}

	frames, err := sampleService.CollectAllFrames(samples)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(frames))

	clusters, err := clusterService.AssignClusters(frames, 500, 50)
	if err != nil {
		panic(err)
	}
	//
	for _, cluster := range clusters {
		err := clusterService.CreateCluster(cluster)
		if err != nil {
			panic(err)
		}
	}

	for _, fm := range frames {
		err := frameService.AssignCluster(fm, clusters)
		if err != nil {
			panic(err)
		}
	}
}
