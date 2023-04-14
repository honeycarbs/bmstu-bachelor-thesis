package repository

import "cluster/internal/entity"

type FrameRepository interface {
	GetOne(sampleHash string, sampleNum int) (entity.Frame, error)
	CountPerSample(sampleHash string) (int, error)
	AssignCluster(clusterID, frameID int) error
}

type SampleRepository interface {
	Get() ([]entity.Sample, error)
}

type ClusterRepository interface {
	Create(entity.Cluster) error
}
