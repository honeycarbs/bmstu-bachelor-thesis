package main

import (
	"fmt"
	"ml/pkg/hmm"
)

func main() {
	m := hmm.HiddenMarkovModel{
		Transitions:             [][]float64{{0.6, 0.4}, {0.3, 0.7}},
		Emissions:               [][]float64{{0.3, 0.7}, {0.8, 0.2}},
		StationaryProbabilities: []float64{0.5, 0.5},
	}

	// Define the observed sequence
	obs := []int{0, 1, 0, 0, 1}

	// Train the HiddenMarkovModel model using the Baum-Welch algorithm
	m.BaumWelch(obs, 1)
	fmt.Println(m.Transitions)
	fmt.Println("-------")
	fmt.Println(m.Emissions)
	//fmt.Println(m.Viterbi(obs))
}
