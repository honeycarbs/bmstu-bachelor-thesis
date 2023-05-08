package hmm

import (
	"math"
)

func (hmm *HiddenMarkovModel) BaumWelch(observationSequence []int, iterations int) {
	var (
		alpha = make([][]float64, len(observationSequence))
		beta  = make([][]float64, len(observationSequence))
		gamma = make([][]float64, len(observationSequence))
		xi    = make([][][]float64, len(observationSequence)-1)
	)

	for i := 0; i < len(observationSequence); i++ {
		alpha[i] = make([]float64, len(hmm.Transitions))
		beta[i] = make([]float64, len(hmm.Transitions))
		gamma[i] = make([]float64, len(hmm.Transitions))
		if i < len(observationSequence)-1 {
			xi[i] = make([][]float64, len(hmm.Transitions))
			for j := 0; j < len(hmm.Transitions); j++ {
				xi[i][j] = make([]float64, len(hmm.Transitions))
			}
		}
	}

	hmm.Emissions = LaplaceSmoothing(hmm.Emissions)
	for it := 0; it < iterations; it++ {
		hmm.ForwardAlgorithm(observationSequence, alpha)
		hmm.BackwardAlgorithm(observationSequence, beta)
		hmm.computeGamma(observationSequence, alpha, beta, gamma)
		hmm.computeXi(observationSequence, alpha, beta, xi)

		// Update the model parameters
		hmm.update(observationSequence, gamma, xi)
	}
}

func (hmm *HiddenMarkovModel) ForwardAlgorithm(observationSequence []int, alpha [][]float64) float64 {
	for i := 0; i < len(hmm.Transitions); i++ {
		alpha[0][i] = math.Log(hmm.StationaryProbabilities[i]) + math.Log(hmm.Emissions[i][observationSequence[0]])
	}

	for t := 1; t < len(observationSequence); t++ {
		for j := 0; j < len(hmm.Transitions); j++ {
			sum := math.Inf(-1)
			for i := 0; i < len(hmm.Transitions); i++ {
				sum = logAdd(sum, alpha[t-1][i]+math.Log(hmm.Transitions[i][j]))
			}
			alpha[t][j] = sum + math.Log(hmm.Emissions[j][observationSequence[t]])
		}
	}

	logLikelihood := math.Inf(-1)
	for i := 0; i < len(hmm.Transitions); i++ {
		logLikelihood = logAdd(logLikelihood, alpha[len(observationSequence)-1][i])
	}
	return logLikelihood
}

func (hmm *HiddenMarkovModel) BackwardAlgorithm(observationSequence []int, beta [][]float64) {
	for i := 0; i < len(hmm.Transitions); i++ {
		beta[len(observationSequence)-1][i] = 0.0
	}

	// Recursion step
	for t := len(observationSequence) - 2; t >= 0; t-- {
		for i := 0; i < len(hmm.Transitions); i++ {
			logSum := math.Inf(-1)
			for j := 0; j < len(hmm.Transitions); j++ {
				logSum = logAdd(logSum, math.Log(hmm.Transitions[i][j])+math.Log(hmm.Emissions[j][observationSequence[t+1]])+beta[t+1][j])
			}
			beta[t][i] = logSum
		}
	}
}

func (hmm *HiddenMarkovModel) computeGamma(obs []int, alpha [][]float64, beta [][]float64, gamma [][]float64) {
	for t := 0; t < len(obs); t++ {
		sum := math.Inf(-1)
		for i := 0; i < len(hmm.Transitions); i++ {
			gamma[t][i] = alpha[t][i] + beta[t][i]
			sum = logAdd(sum, gamma[t][i])
		}

		for i := 0; i < len(hmm.Transitions); i++ {
			gamma[t][i] -= sum
			gamma[t][i] = math.Exp(gamma[t][i])
		}
	}
}

func (hmm *HiddenMarkovModel) computeXi(obs []int, alpha [][]float64, beta [][]float64, xi [][][]float64) {
	for t := 0; t < len(obs)-1; t++ {
		for i := 0; i < len(hmm.Transitions); i++ {
			for j := 0; j < len(hmm.Transitions); j++ {
				xi[t][i][j] = alpha[t][i] + math.Log(hmm.Transitions[i][j]) + math.Log(hmm.Emissions[j][obs[t+1]]) + beta[t+1][j]
			}
		}
	}

	for t := 0; t < len(obs)-1; t++ {
		maxVal := math.Inf(-1)
		for i := 0; i < len(hmm.Transitions); i++ {
			for j := 0; j < len(hmm.Transitions); j++ {
				if xi[t][i][j] > maxVal {
					maxVal = xi[t][i][j]
				}
			}
		}

		sum := 0.0
		for i := 0; i < len(hmm.Transitions); i++ {
			for j := 0; j < len(hmm.Transitions); j++ {
				xi[t][i][j] = math.Exp(xi[t][i][j] - maxVal)
				sum += xi[t][i][j]
			}
		}

		for i := 0; i < len(hmm.Transitions); i++ {
			for j := 0; j < len(hmm.Transitions); j++ {
				xi[t][i][j] /= sum
			}
		}
	}
}

func (hmm *HiddenMarkovModel) update(obs []int, gamma [][]float64, xi [][][]float64) {
	// Update initial state probabilities
	for i := 0; i < len(hmm.StationaryProbabilities); i++ {
		hmm.StationaryProbabilities[i] = gamma[0][i]
	}

	// Update transition probabilities
	for i := 0; i < len(hmm.Transitions); i++ {
		for j := 0; j < len(hmm.Transitions); j++ {
			sumXi := 0.0
			sumGamma := 0.0
			for t := 0; t < len(obs)-1; t++ {
				sumXi += xi[t][i][j]
				sumGamma += gamma[t][i]
			}
			hmm.Transitions[i][j] = sumXi / sumGamma
		}
	}

	// Update emission probabilities
	for i := 0; i < len(hmm.Emissions); i++ {
		for j := 0; j < len(hmm.Emissions[0]); j++ {
			sumGamma := 0.0
			sumGammaObs := 0.0
			for t := 0; t < len(obs); t++ {
				if obs[t] == j {
					sumGammaObs += gamma[t][i]
				}
				sumGamma += gamma[t][i]
			}
			hmm.Emissions[i][j] = sumGammaObs / sumGamma
		}
	}
}

func LaplaceSmoothing(emissionMatrix [][]float64) [][]float64 {
	rows := len(emissionMatrix)
	cols := len(emissionMatrix[0])
	smoothedMatrix := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		smoothedMatrix[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			smoothedMatrix[i][j] = (emissionMatrix[i][j] + 1) / (sum(emissionMatrix[i]) + float64(cols))
		}
	}
	return smoothedMatrix
}
