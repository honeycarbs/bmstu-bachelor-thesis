package service

import (
	"ml/internal/entity"
	"ml/internal/repository/postgres"
)

type SampleService struct {
	sampleRepo *postgres.SamplePostgres
	frameRepo  *postgres.FramePostgres
}

func NewSampleService(sampleRepo *postgres.SamplePostgres, frameRepo *postgres.FramePostgres) *SampleService {
	return &SampleService{sampleRepo: sampleRepo, frameRepo: frameRepo}
}

func (s *SampleService) GetByLabel(label entity.Label) ([]entity.Sample, error) {
	samples, err := s.sampleRepo.GetByLabel(label)
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

func (s *SampleService) ConstructObservationSequence(sample entity.Sample) []int {
	observations := make([]int, len(sample.Frames))

	for i := 0; i < len(observations); i++ {
		observations[i] = sample.Frames[i].ClusterIndex - 1
	}

	return observations
}

//func (s *SampleService) ComputeEmpiricalClusterProbabilities(samples []entity.Sample, nClusters int) [][]float64 {
//	// initialize the emission matrix with zeros
//	emissionMatrix := make([][]float64, 1)
//	for i := range emissionMatrix {
//		emissionMatrix[i] = make([]float64, nClusters)
//	}
//
//	// compute the empirical frequencies for each hidden state and observed state
//	for _, sample := range samples {
//		for _, frame := range sample.Frames {
//			// increment the count for the observed state in the appropriate hidden state
//			emissionMatrix[0][frame.ClusterIndex-1]++
//		}
//	}
//
//	// normalize the counts to get probabilities
//	for i := range emissionMatrix {
//		sum := 0.0
//		for j := range emissionMatrix[i] {
//			sum += emissionMatrix[i][j]
//		}
//		for j := range emissionMatrix[i] {
//			emissionMatrix[i][j] /= sum
//		}
//	}
//
//	return emissionMatrix
//}
