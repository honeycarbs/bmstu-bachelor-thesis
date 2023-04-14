package repository

import "ml/internal/entity"

type ClusterRepository interface {
	GetByFrame(frameUUID string) ([]entity.Cluster, error)
}

type FrameRepository interface {
	getBySample(sampleUUID string) []entity.Frame
}

type SampleRepository interface {
	getByLabel(label string) ([]entity.Sample, error)
}
