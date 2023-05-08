package repository

import "ml/internal/entity"

type ClusterRepository interface {
	GetByFrame(frameUUID string) ([]entity.Cluster, error)
}

type FrameRepository interface {
	GetBySample(sampleUUID string) []entity.Frame
}

type SampleRepository interface {
	GetByLabelTrain(label entity.Label) ([]entity.Sample, error)
}
