package repository

import "cluster/internal/entity"

type FrameRepository interface {
	AssignCluster(clusterID, frameID int) error
	GetAll() ([]entity.Frame, error)
	GetByAudio(path string) ([]entity.Frame, error)
}

type ClusterRepository interface {
	Create(entity.Cluster) error
	Get() ([]entity.Cluster, error)
}
