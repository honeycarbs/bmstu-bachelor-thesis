package hmm

import "math/rand"

type HiddenMarkovModel struct {
	Transitions             [][]float64 `json:"transitions"`
	Emissions               [][]float64 `json:"emissions"`
	StationaryProbabilities []float64   `json:"stationary_probabilities"`
}

func New(nStates, nObs int) *HiddenMarkovModel {
	seed := rand.New(rand.NewSource(0))

	emitProb := allocateRandomMatrix(nStates, nObs, seed)
	transProb := allocateRandomMatrix(nStates, nStates, seed)

	initProb := make([]float64, nStates)
	sum := 0.
	for i := range initProb {
		initProb[i] = seed.Float64()
		sum += initProb[i]
	}

	for i := range initProb {
		initProb[i] /= sum
	}
	return &HiddenMarkovModel{
		Transitions:             transProb,
		Emissions:               emitProb,
		StationaryProbabilities: initProb,
	}
}

func allocateRandomMatrix(n, m int, seed *rand.Rand) [][]float64 {
	matrix := make([][]float64, n)
	for i := range matrix {
		matrix[i] = make([]float64, m)
		sum := 0.0
		for j := range matrix[i] {
			matrix[i][j] = seed.Float64()
			sum += matrix[i][j]
		}

		for j := range matrix[i] {
			matrix[i][j] /= sum
		}
	}

	return matrix
}
