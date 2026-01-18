
package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	GeminiModel   string  `json:"gemini_model"`
	Temperature   float32 `json:"temperature"`
	MaxTokens     int     `json:"max_tokens"`
	MemoryType    string  `json:"memory_type"`
	DBConnectionString string `json:"db_connection_string"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
