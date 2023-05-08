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

func (s *FrameService) GetAllBySample(sampleUUID string) ([]entity.Frame, error) {
	count, err := s.countFramesPerSample(sampleUUID)
	if err != nil {
		return nil, err
	}

	frames := make([]entity.Frame, count)
	for i := 1; i <= count; i++ {
		frame, err := s.repo.GetOne(sampleUUID, i)
		if err != nil {
			return nil, err
		}
		frames[i-1] = frame
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

func (s *FrameService) countFramesPerSample(sampleHash string) (int, error) {
	count, err := s.repo.CountPerSample(sampleHash)
	if err != nil {
		return 0, err
	}
	return count, nil
}
