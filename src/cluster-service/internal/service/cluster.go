package service

import (
	"cluster/internal/entity"
	"cluster/internal/repository/postgres"
	"cluster/pkg/kmeans"
	"cluster/pkg/logging"
	"github.com/google/uuid"
)

type ClusterService struct {
	repo   *postgres.ClusterPostgres
	logger logging.Logger
}

func NewClusterService(repo *postgres.ClusterPostgres, logger logging.Logger) *ClusterService {
	return &ClusterService{repo: repo, logger: logger}
}

func (s *ClusterService) AssignClusters(frames []entity.Frame, nclusters, maxRounds int) ([]entity.Cluster, error) {
	nodes := s.collectFramesData(frames)
	centroidsCoords, err := kmeans.KMeans(nodes, nclusters, maxRounds)
	if err != nil {
		return nil, err
	}

	centroids := s.constructCentroidsData(centroidsCoords)
	return s.constructClusterData(centroids), nil
}

func (s *ClusterService) CreateCluster(cluster entity.Cluster) error {
	return s.repo.Create(cluster)
}

func (s *ClusterService) collectFramesData(frames []entity.Frame) []kmeans.Node {
	nodes := make([]kmeans.Node, len(frames))

	for i, fm := range frames {
		nodes[i] = fm.MFCCs
	}

	return nodes
}

func (s *ClusterService) constructCentroidsData(centroidsCoords []kmeans.Node) []entity.Centroid {
	centroids := make([]entity.Centroid, len(centroidsCoords))

	for i, centroid := range centroidsCoords {
		centroids[i] = entity.Centroid{
			ID:    uuid.New().String(),
			Value: centroid,
		}
	}

	return centroids
}

func (s *ClusterService) constructClusterData(centroids []entity.Centroid) []entity.Cluster {
	clusters := make([]entity.Cluster, len(centroids))
	for i, centroid := range centroids {
		clusters[i] = entity.Cluster{
			ID:       uuid.New().String(),
			Index:    i + 1,
			Centroid: centroid,
		}
	}
	return clusters
}
