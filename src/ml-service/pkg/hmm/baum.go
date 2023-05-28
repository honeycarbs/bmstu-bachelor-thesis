package hmm

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

	for it := 0; it < iterations; it++ {
		hmm.Forward(observationSequence, alpha)
		hmm.Backward(observationSequence, beta)

		hmm.computeGammaProbs(observationSequence, alpha, beta, gamma)
		hmm.computeXiProbs(observationSequence, alpha, beta, xi)
		
		hmm.update(observationSequence, gamma, xi)
	}

	hmm.laplaceSmoothing(observationSequence)
}

func (hmm *HiddenMarkovModel) computeGammaProbs(obs []int, alpha [][]float64, beta [][]float64, gamma [][]float64) {
	for t := 0; t < len(obs); t++ {
		sum := 0.0
		for i := 0; i < len(hmm.Transitions); i++ {
			gamma[t][i] = alpha[t][i] * beta[t][i]
			sum += gamma[t][i]
		}

		for i := 0; i < len(hmm.Transitions); i++ {
			gamma[t][i] /= sum
		}
	}
}

func (hmm *HiddenMarkovModel) computeXiProbs(obs []int, alpha [][]float64, beta [][]float64, xi [][][]float64) {
	for t := 0; t < len(obs)-1; t++ {
		sum := 0.0
		for i := 0; i < len(hmm.Transitions); i++ {
			for j := 0; j < len(hmm.Transitions); j++ {
				xi[t][i][j] = alpha[t][i] * hmm.Transitions[i][j] * hmm.Emissions[j][obs[t+1]] * beta[t+1][j]
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
	for i := 0; i < len(hmm.StationaryProbabilities); i++ {
		hmm.StationaryProbabilities[i] = gamma[0][i]
	}

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
