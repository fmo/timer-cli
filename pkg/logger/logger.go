// Package logger log wrapper
package logger

import (
	"log"
	"os"
)

func New() (*log.Logger, error) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	logger := log.New(file, "log: ", log.Lshortfile)
	return logger, nil
}
