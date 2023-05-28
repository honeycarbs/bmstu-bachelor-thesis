package hmm

const (
	SMOOTHING_FACTOR float64 = 0.5
)

func (hmm *HiddenMarkovModel) laplaceSmoothing(observations []int) {
	nStates := len(hmm.Emissions)
	nObs := len(hmm.Emissions[0])

	emissions := make([][]float64, nStates)
	for i := range emissions {
		emissions[i] = make([]float64, nObs)
	}

	counts := make([]int, nObs)
	for i := 0; i < len(observations); i++ {
		state := observations[i]
		counts[state]++
	}

	total := 0
	for _, count := range counts {
		total += count
	}

	for i := 0; i < nObs; i++ {
		smoothedProbability := (float64(counts[i]) + SMOOTHING_FACTOR) / (float64(total) + SMOOTHING_FACTOR*float64(nObs))
		emissions[0][i] = smoothedProbability
	}
	hmm.Emissions = emissions
}
