package service

import (
	"fmt"
	"math"
	"ml/internal/config"
	"ml/internal/entity"
	"ml/internal/pkg/heatmap"
	"ml/pkg/hmm"
	"ml/pkg/logging"
	"time"
)

var (
	labels = []entity.Label{entity.NEUTRAL, entity.ANGRY, entity.POSITIVE, entity.SAD}
	//labels = []entity.Label{entity.ANGRY}
)

type MLService struct {
	logger        logging.Logger
	sampleService *SampleService
}

func NewMLService(logger logging.Logger, service *SampleService) *MLService {
	return &MLService{logger: logger, sampleService: service}
}

func (s *MLService) Train() error {
	start := time.Now()
	for _, label := range labels {
		s.logger.Infof("Started training model for (%v)...", label)
		samples, err := s.sampleService.GetByLabelTrain(label)
		if err != nil {
			return err
		}

		model := hmm.New(1, config.GetConfig().ClusterAmount)
		for i, sm := range samples {
			obs := s.sampleService.ConstructObservationSequence(sm)
			model.BaumWelch(obs, 100)
			s.logger.Infof("trained %v sample out of %v for model %v, NaN encountered: %v",
				i, len(samples), label, hasNaN(model.Emissions))
		}

		err = model.SaveJSON(fmt.Sprintf(config.GetConfig().FilePath, label))
		if err != nil {
			return err
		}
	}
	s.logger.Infof("Trained all models, took %v ", time.Since(start).Milliseconds())

	return nil
}

func (s *MLService) Test() ([]entity.Label, []entity.Label, error) {
	actual := make([]entity.Label, 0)
	predicted := make([]entity.Label, 0)

	for _, label := range labels {
		s.logger.Infof("Started testing for label %v", label)
		samples, err := s.sampleService.GetByLabelTest(label)
		if err != nil {
			panic(err)
		}

		models := make([]hmm.HiddenMarkovModel, len(labels))
		for i, label := range labels {
			models[i], err = hmm.LoadJSON(fmt.Sprintf(config.GetConfig().FilePath, label))
			if err != nil {
				return nil, nil, fmt.Errorf("models are not trained, training required")
			}
		}

		for _, sample := range samples {
			observationSequence := s.sampleService.ConstructObservationSequence(sample)
			index := hmm.FindBestFittedModel(observationSequence, models)

			actual = append(actual, label)
			predicted = append(predicted, labels[index])

			if label == labels[index] && label != entity.NEUTRAL {
				fmt.Println(sample.AudioPath, label)
			}
		}
	}
	fmt.Println(actual)
	fmt.Println(predicted)

	return actual, predicted, nil
}

func (s *MLService) GetHeatmap(actual, predicted []entity.Label) {
	heatmap.GetHeatmap(actual, predicted)
	s.logger.Infof("heatmap obtained")
}

func (s *MLService) Recognize(path string) (entity.Label, error) {
	sample, err := s.sampleService.GetByPath(path)
	if err != nil {
		return "", err
	}

	models := make([]hmm.HiddenMarkovModel, len(labels))
	for i, label := range labels {
		models[i], err = hmm.LoadJSON(fmt.Sprintf(config.GetConfig().FilePath, label))
		if err != nil {
			return "", fmt.Errorf("models are not trained, training required")
		}
	}

	observationSequence := s.sampleService.ConstructObservationSequence(sample)
	index := hmm.FindBestFittedModel(observationSequence, models)

	return labels[index], nil
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

// TODO: remove this bitch
func hasNaN(matrix [][]float64) bool {
	for _, row := range matrix {
		for _, val := range row {
			if math.IsNaN(val) {
				return true
			}
		}
	}
	return false
}
