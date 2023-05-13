package hmm

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"ml/pkg/logging"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestHiddenMarkovModel_BaumWelch(t *testing.T) {
	testCases := []struct {
		name  string
		fEmit string
		fObs  string
		fExp  string
		dim   int
	}{
		{
			name:  "STATES DIFFER",
			fEmit: "etc/test-cases/hmm-diff-emit",
			fExp:  "etc/test-cases/hmm-diff-expect",
			fObs:  "etc/test-cases/hmm-diff-obs",
			dim:   50,
		},
		{
			name:  "SAME STATE",
			fEmit: "etc/test-cases/hmm-same-emit",
			fExp:  "etc/test-cases/hmm-same-expect",
			fObs:  "etc/test-cases/hmm-same-obs",
			dim:   50,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := runBaumWelchTest(testCase.fEmit, testCase.fObs, testCase.fExp, testCase.dim)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestHiddenMarkovModel_BaumWelchBigMatrix(t *testing.T) {
	dim := 2500
	model := New(1, dim)
	obs, err := readArrayFromFile("etc/test-cases/hmm-big-obs")
	if err != nil {
		t.Fatal(err)
	}

	model.BaumWelch(obs, 100)

	if hasNaN(model.Emissions) {
		t.Fatal(fmt.Errorf("matrix has NaN values"))
	}
}

// // // // // // // // // // // helpers // // // // // // // // // // //
func runBaumWelchTest(fEmit, fObs, fExp string, dim int) error {
	emit, err := readMatrixFromFile(fEmit, 1, dim)
	if err != nil {
		return err
	}

	hmm := HiddenMarkovModel{
		Transitions: [][]float64{
			{1},
		},
		Emissions:               emit,
		StationaryProbabilities: []float64{1},
	}

	obs, err := readArrayFromFile(fObs)
	if err != nil {
		return err
	}

	hmm.BaumWelch(obs, 1)
	expectedEmissions, err := readMatrixFromFile(fExp, 1, dim)
	if err != nil {
		return err
	}

	equal := matricesEqual(hmm.Emissions, expectedEmissions, 0.005)
	if !equal {
		return errors.New("emission matrix differs from expected")
	}

	return nil
}

func matricesEqual(matrix1 [][]float64, matrix2 [][]float64, epsilon float64) bool {
	if len(matrix1) != len(matrix2) || len(matrix1[0]) != len(matrix2[0]) {
		return false
	}

	for i := 0; i < len(matrix1); i++ {
		for j := 0; j < len(matrix1[0]); j++ {
			diff := matrix1[i][j] - matrix2[i][j]
			if diff < -epsilon || diff > epsilon {
				logging.Init()
				logging.GetLogger().Infof("%v %v %v", matrix1[i][j], matrix2[i][j], matrix1[i][j]-matrix2[i][j])
				return false
			}
		}
	}

	return true
}

func readMatrixFromFile(filePath string, rows, cols int) ([][]float64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	matrix := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]float64, cols)
		if !scanner.Scan() {
			return nil, fmt.Errorf("%v: insufficient data in file", filePath)
		}
		line := scanner.Text()
		nums := strings.Fields(line)
		if len(nums) != cols {
			return nil, fmt.Errorf("%v: incorrect number of columns in line %d: %v", filePath, i+1, len(nums))
		}
		for j, numStr := range nums {
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return nil, fmt.Errorf("%v: failed to parse number in line %d: %v", filePath, i+1, err)
			}
			matrix[i][j] = num
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return matrix, nil
}

func readArrayFromFile(filePath string) ([]int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var array []int

	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Fields(line)
		for _, numStr := range nums {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return nil, fmt.Errorf("failed to parse number: %v", err)
			}
			array = append(array, num)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return array, nil
}

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
