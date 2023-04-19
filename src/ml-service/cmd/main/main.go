package main

import (
	"fmt"
	"ml/internal/entity"
	"ml/internal/repository/postgres"
	"ml/internal/service"
	"ml/pkg/hmm"
	"ml/pkg/psqlcli"
)

const (
	filepathTemplate = "_meta/model-%s.json"
	LabelNum         = 4
)

// TODO распознавание: вытащить из репозитория по разметке -> запустить для каждой  модели алгоритм прямого хода -> найти максимум
func main() {
	cli, err := psqlcli.New("localhost", "5432", "admin", "admin", "thesis", "disable")
	if err != nil {
		panic(err)
	}

	sampleRepo := postgres.NewSamplePostgres(cli)
	frameRepo := postgres.NewFramePostgres(cli)

	sampleService := service.NewSampleService(sampleRepo, frameRepo)

	err = test(sampleService)

	//err = train(entity.NEUTRAL, sampleService)
	if err != nil {
		panic(err)
	}
	//
	//err = train(entity.POSITIVE, sampleService)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = train(entity.SAD, sampleService)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = train(entity.ANGRY, sampleService)
	//if err != nil {
	//	panic(err)
	//}
}

//func train(label entity.Label, sampleService *service.SampleService) error {
//	samples, err := sampleService.GetByLabel(label)
//	if err != nil {
//		panic(err)
//	}
//	model := hmm.New(1, 500)
//	for _, s := range samples {
//		obs := sampleService.ConstructObservationSequence(s)
//		model.BaumWelch(obs, 10)
//		fmt.Println("--- emissions AFTER BAUM WELCH -----")
//		fmt.Println(model.Emissions)
//	}
//	err = model.SaveJSON(fmt.Sprintf(filepathTemplate, label))
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func test(sampleService *service.SampleService) error {
	labels := []entity.Label{entity.NEUTRAL, entity.POSITIVE, entity.ANGRY, entity.SAD}
	actual := make([]entity.Label, 0)
	predicted := make([]entity.Label, 0)

	for _, label := range labels {
		samples, err := sampleService.GetByLabel(label)
		if err != nil {
			panic(err)
		}

		models := make([]hmm.HiddenMarkovModel, len(labels))
		for i, label := range labels {
			models[i], err = hmm.LoadJSON(fmt.Sprintf(filepathTemplate, label))
			if err != nil {
				return err
			}
		}

		for _, sample := range samples {
			observationSequence := sampleService.ConstructObservationSequence(sample)
			likelihoods := make([]float64, len(models))
			actual = append(actual, label)

			fmt.Println(len(models))
			for i, model := range models {
				alpha := make([][]float64, len(observationSequence))
				for i := 0; i < len(observationSequence); i++ {
					alpha[i] = make([]float64, len(model.Transitions))
				}

				model.Emissions = hmm.LaplaceSmoothing(model.Emissions)
				likelihood := model.ForwardAlgorithm(observationSequence, alpha)
				likelihoods[i] = likelihood
			}
			//fmt.Println(likelihoods)

			predictedIndex, _ := findMax(likelihoods)
			predicted = append(predicted, labels[predictedIndex])
		}

		//modelNetural, err := hmm.LoadJSON(fmt.Sprintf(filepathTemplate, entity.NEUTRAL))
		//modelPositive, err := hmm.LoadJSON(fmt.Sprintf(filepathTemplate, entity.POSITIVE))
		//modelSadAngry, err := hmm.LoadJSON(fmt.Sprintf(filepathTemplate, entity.ANGRY))
		//modelSad, err := hmm.LoadJSON(fmt.Sprintf(filepathTemplate, entity.SAD))

		//likelihoods := make([]float64, len(labels))
		//for i := 0; i < len(likelihoods); i++ {
		//	likelihood :=
		//}
	}

	confusionMatrix(actual, predicted)

	//for i := 0; i < len(actual); i++ {
	//	fmt.Println(actual[i], predicted[i])
	//}
	//samples, err := sampleService.GetByLabel(label)
	//if err != nil {
	//	panic(err)
	//}
	//
	//likelihoods := make([]float64, LabelNum)
	//
	//modelNetural := hmm.New(1, 500)
	//fmt.Println(actual)
	//fmt.Println(predicted)
	return nil
}

func confusionMatrix(actual, predicted []entity.Label) {
	// Initialize the confusion matrix with zeros
	cm := make(map[entity.Label]map[entity.Label]int)
	for _, actualClass := range []entity.Label{entity.NEUTRAL, entity.POSITIVE, entity.ANGRY, entity.SAD} {
		cm[actualClass] = make(map[entity.Label]int)
		for _, predictedClass := range []entity.Label{entity.NEUTRAL, entity.POSITIVE, entity.ANGRY, entity.SAD} {
			cm[actualClass][predictedClass] = 0
		}
	}

	// Fill the confusion matrix by iterating over each actual and predicted label
	for i := 0; i < len(actual); i++ {
		cm[actual[i]][predicted[i]]++
	}

	// Print the confusion matrix
	fmt.Println("Confusion Matrix:")
	fmt.Printf("%10s%10s%10s%10s%10s\n", "", "Predicted", "", "", "")
	fmt.Printf("%10s%10s%10s%10s%10s\n", "", "Neutral", "Positive", "Angry", "Sad")
	fmt.Printf("%10s%10s%10s%10s%10s\n", "Actual", "--------", "--------", "--------", "--------")
	for _, actualClass := range []entity.Label{entity.NEUTRAL, entity.POSITIVE, entity.ANGRY, entity.SAD} {
		fmt.Printf("%10s%10d%10d%10d%10d\n", actualClass, cm[actualClass][entity.NEUTRAL],
			cm[actualClass][entity.POSITIVE], cm[actualClass][entity.ANGRY],
			cm[actualClass][entity.SAD])
	}
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

func getIndex(label entity.Label) int {
	switch label {
	case entity.SAD:
		return 0
	case entity.ANGRY:
		return 1
	case entity.POSITIVE:
		return 2
	default:
		return -1 // обработка ошибок
	}
}
