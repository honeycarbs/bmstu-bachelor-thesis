package service

import (
	"fmt"
	"ml/internal/config"
	"ml/internal/entity"
	"ml/internal/pkg/heatmap"
	"ml/pkg/hmm"
	"ml/pkg/logging"
)

type MLService struct {
	logger        logging.Logger
	sampleService *SampleService
}

func NewMLService(logger logging.Logger, service *SampleService) *MLService {
	return &MLService{logger: logger, sampleService: service}
}

func (s *MLService) Train() error {
	labels := []entity.Label{entity.NEUTRAL, entity.POSITIVE, entity.ANGRY, entity.SAD}

	for _, label := range labels {
		samples, err := s.sampleService.GetByLabelTrain(label)
		if err != nil {
			return err
		}
		model := hmm.New(1, config.GetConfig().ClusterAmount)
		for _, sm := range samples {
			obs := s.sampleService.ConstructObservationSequence(sm)
			model.BaumWelch(obs, 100)
		}

		err = model.SaveJSON(fmt.Sprintf(config.GetConfig().FilePath, label))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MLService) Test() ([]entity.Label, []entity.Label, error) {
	labels := []entity.Label{entity.NEUTRAL, entity.POSITIVE, entity.ANGRY, entity.SAD}
	actual := make([]entity.Label, 0)
	predicted := make([]entity.Label, 0)

	for _, label := range labels {
		samples, err := s.sampleService.GetByLabelTest(label)
		if err != nil {
			panic(err)
		}

		models := make([]hmm.HiddenMarkovModel, len(labels))
		for i, label := range labels {
			models[i], err = hmm.LoadJSON(fmt.Sprintf(config.GetConfig().FilePath, label))
			if err != nil {
				return nil, nil, err
			}
		}

		for i, sample := range samples {
			observationSequence := s.sampleService.ConstructObservationSequence(sample)
			likelihoods := make([]float64, len(models))
			actual = append(actual, label)

			s.logger.Infof("Tested (%v) samples out of (%v)...", i, len(samples))

			//fmt.Println(len(models))
			for i, model := range models {
				alpha := make([][]float64, len(observationSequence))
				for i := 0; i < len(observationSequence); i++ {
					alpha[i] = make([]float64, len(model.Transitions))
				}

				model.Emissions = hmm.LaplaceSmoothing(model.Emissions)
				likelihood := model.ForwardAlgorithm(observationSequence, alpha)
				likelihoods[i] = likelihood
			}

			predictedIndex, _ := findMax(likelihoods)
			predicted = append(predicted, labels[predictedIndex])
		}
	}

	s.logger.Printf("actual: %v", actual)
	s.logger.Printf("predicted: %v", predicted)
	return actual, predicted, nil
}

func (s *MLService) GetHeatmap(actual, predicted []entity.Label) {
	heatmap.GetHeatmap(actual, predicted)
	s.logger.Infof("heatmap obtained")
}

func findMax(slice []float64) (int, float64) {
	max := slice[0]
	index := 0
	for i, value := range slice {
		if value > max {
			max = value
			index = i
		}
	}
	return index, max
}
