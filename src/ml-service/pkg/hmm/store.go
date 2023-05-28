package hmm

import (
	"encoding/json"
	"os"
)

func (hmm *HiddenMarkovModel) SaveJSON(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(hmm); err != nil {
		return err
	}

	return nil
}

func LoadJSON(filename string) (HiddenMarkovModel, error) {
	file, err := os.Open(filename)
	if err != nil {
		return HiddenMarkovModel{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var hmm HiddenMarkovModel
	if err := decoder.Decode(&hmm); err != nil {
		return HiddenMarkovModel{}, err
	}

	return hmm, nil
}
