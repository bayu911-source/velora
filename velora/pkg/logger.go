
package utils

import (
	"log"
	"os"
)

func NewLogger(logPath string) (*log.Logger, error) {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	return log.New(file, "", log.LstdFlags), nil
}
