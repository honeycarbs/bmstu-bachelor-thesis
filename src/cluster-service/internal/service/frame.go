package service

import (
	"cluster/internal/entity"
	"cluster/internal/repository/postgres"
	"cluster/pkg/kmeans"
)

type FrameService struct {
	repo *postgres.FramePostgres
}

func NewFrameService(repo *postgres.FramePostgres) *FrameService {
	return &FrameService{repo: repo}
}

func (s *FrameService) GetAll() ([]entity.Frame, error) {
	frames, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return frames, nil
}

func (s *FrameService) AssignCluster(frame entity.Frame, clusters []entity.Cluster) error {
	centroids := make([]kmeans.Node, len(clusters))
	for i, cluster := range clusters {
		centroids[i] = cluster.Centroid.Value
	}
	nearestIndex := kmeans.Nearest(frame.MFCCs, centroids)
	nearest := clusters[nearestIndex]

	return s.repo.AssignCluster(nearest.ID, frame.ID)
}
