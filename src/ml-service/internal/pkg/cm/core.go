package cm

import "ml/internal/entity"

type ConfusionMatrix struct {
	Values map[entity.Label]map[entity.Label]float64
	labels []entity.Label
}

func NewConfusionMatrix(labels []entity.Label, actual, predicted []entity.Label) *ConfusionMatrix {
	confusionMatrix := make(map[entity.Label]map[entity.Label]float64)
	for _, label := range labels {
		confusionMatrix[label] = make(map[entity.Label]float64)
		for _, label2 := range labels {
			confusionMatrix[label][label2] = 0.
		}
	}

	for i := 0; i < len(actual); i++ {
		confusionMatrix[actual[i]][predicted[i]]++
	}

	return &ConfusionMatrix{labels: labels, Values: confusionMatrix}
}

func (m *ConfusionMatrix) Normalize() {
	// Compute the number of classes
	numClasses := len(m.labels)

	// Compute the row sums and the total sum of the confusion matrix
	rowSums := make([]float64, numClasses)
	totalSum := 0.0
	for i, actualLabel := range m.labels {
		for _, predictedLabel := range m.labels {
			rowSums[i] += m.Values[actualLabel][predictedLabel]
			totalSum += m.Values[actualLabel][predictedLabel]
		}
	}

	// Compute the normalized confusion matrix
	normConfusionMatrix := make(map[entity.Label]map[entity.Label]float64)

	for i, actualLabel := range m.labels {
		normConfusionMatrix[actualLabel] = make(map[entity.Label]float64)
		for _, predictedLabel := range m.labels {
			normConfusionMatrix[actualLabel][predictedLabel] = float64(m.Values[actualLabel][predictedLabel]) / rowSums[i]
		}
	}
	m.Values = normConfusionMatrix
}
