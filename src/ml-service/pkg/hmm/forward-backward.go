package hmm

import (
	"gonum.org/v1/gonum/floats"
)

func (hmm *HiddenMarkovModel) Forward(observationSequence []int, alpha [][]float64) {
	for i := 0; i < len(hmm.Transitions); i++ {
		alpha[0][i] = hmm.StationaryProbabilities[i] * hmm.Emissions[i][observationSequence[0]]
	}

	for t := 1; t < len(observationSequence); t++ {
		for j := 0; j < len(hmm.Transitions); j++ {
			sum := 0.0
			for i := 0; i < len(hmm.Transitions); i++ {
				sum += alpha[t-1][i] * hmm.Transitions[i][j]
			}
			alpha[t][j] = sum * hmm.Emissions[j][observationSequence[t]]
		}
	}
}

func (hmm *HiddenMarkovModel) Backward(observationSequence []int, beta [][]float64) {
	for i := 0; i < len(hmm.Transitions); i++ {
		beta[len(observationSequence)-1][i] = 1.0
	}

	for t := len(observationSequence) - 2; t >= 0; t-- {
		for i := 0; i < len(hmm.Transitions); i++ {
			sum := 0.0
			for j := 0; j < len(hmm.Transitions); j++ {
				sum += hmm.Transitions[i][j] * hmm.Emissions[j][observationSequence[t+1]] * beta[t+1][j]
			}
			beta[t][i] = sum
		}
	}
}

func FindBestFittedModel(observationSequence []int, models []HiddenMarkovModel) int {
	probabilities := make([]float64, len(models))

	for i, model := range models {
		alpha := make([][]float64, len(observationSequence))
		for i := 0; i < len(observationSequence); i++ {
			alpha[i] = make([]float64, len(model.Transitions))
		}

		model.Forward(observationSequence, alpha)
		probability := 0.0
		for i := 0; i < len(models[0].Transitions); i++ {
			probability += alpha[len(observationSequence)-1][i]
		}
		probabilities[i] = probability
	}

	max := floats.MaxIdx(probabilities)

	return max
}
