
package workflow

import (
	"encoding/json"
	"os"
)

type Step struct {
	Agent string `json:"agent"`
	Input string `json:"input"`
}

type Pipeline struct {
	Name  string `json:"name"`
	Steps []Step `json:"steps"`
}

func LoadPipeline(path string) (*Pipeline, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	pipeline := &Pipeline{}
	err = decoder.Decode(pipeline)
	if err != nil {
		return nil, err
	}

	return pipeline, nil
}
