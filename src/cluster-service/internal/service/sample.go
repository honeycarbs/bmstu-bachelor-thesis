package service

import (
	"cluster/internal/entity"
	"cluster/internal/repository/postgres"
)

type SampleService struct {
	repo *postgres.SamplePostgres
}

func NewSampleService(repo *postgres.SamplePostgres) *SampleService {
	return &SampleService{repo: repo}
}

func (s *SampleService) GetAll() ([]entity.Sample, error) {
	var samples []entity.Sample

	samples, err := s.repo.Get()
	if err != nil {
		return nil, err
	}

	return samples, nil
}

func (s *SampleService) CollectAllFrames(samples []entity.Sample) ([]entity.Frame, error) {
	framesLength := len(samples) * len(samples[0].Frames)

	frames := make([]entity.Frame, framesLength)
	index := 0
	for i := 0; i < len(samples); i++ {
		for j := 0; j < len(samples[i].Frames); j++ {
			frames[index] = samples[i].Frames[j]
			index++
		}
	}

	return frames, nil
}
