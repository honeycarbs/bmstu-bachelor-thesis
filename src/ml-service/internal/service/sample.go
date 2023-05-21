package service

import (
	"ml/internal/entity"
	"ml/internal/repository/postgres"
	"ml/pkg/logging"
)

type SampleService struct {
	logger     logging.Logger
	sampleRepo *postgres.SamplePostgres
	frameRepo  *postgres.FramePostgres
}

func NewSampleService(sampleRepo *postgres.SamplePostgres, frameRepo *postgres.FramePostgres, logger logging.Logger) *SampleService {
	return &SampleService{logger: logger, sampleRepo: sampleRepo, frameRepo: frameRepo}
}

func (s *SampleService) GetByLabelTrain(label entity.Label) ([]entity.Sample, error) {
	samples, err := s.sampleRepo.GetByLabelTrain(label)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(samples); i++ {
		frames, err := s.frameRepo.GetBySample(samples[i].ID)
		if err != nil {
			return nil, err
		}
		samples[i].Frames = frames
	}

	return samples, nil
}

func (s *SampleService) GetByLabelTest(label entity.Label) ([]entity.Sample, error) {
	samples, err := s.sampleRepo.GetByLabelTest(label)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(samples); i++ {
		frames, err := s.frameRepo.GetBySample(samples[i].ID)
		if err != nil {
			return nil, err
		}
		samples[i].Frames = frames
	}

	return samples, nil
}

func (s *SampleService) GetByPath(path string) (entity.Sample, error) {
	sample, err := s.sampleRepo.GetByPath(path)
	if err != nil {
		return entity.Sample{}, err
	}

	frames, err := s.frameRepo.GetBySample(sample.ID)
	if err != nil {
		return entity.Sample{}, err
	}
	sample.Frames = frames

	return sample, nil
}

func (s *SampleService) ConstructObservationSequence(sample entity.Sample) []int {
	observations := make([]int, len(sample.Frames))

	for i := 0; i < len(observations); i++ {
		observations[i] = sample.Frames[i].ClusterIndex
	}

	return observations
}
