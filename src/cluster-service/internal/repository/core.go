package repository

import "cluster/internal/entity"

type FrameRepository interface {
	AssignCluster(clusterID, frameID int) error
	GetAll() ([]entity.Frame, error)
}

type ClusterRepository interface {
	Create(entity.Cluster) error
}
