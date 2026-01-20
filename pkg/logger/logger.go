// Package logger log wrapper
package logger

import (
	"log"
	"os"
)

type Logger interface {
	Fatal(v ...any)
	Fatalf(format string, v ...any)
}

type LoggerImpl struct {
	l *log.Logger
}

func New() (*LoggerImpl, error) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	logger := log.New(file, "log: ", log.Lshortfile)
	return &LoggerImpl{l: logger}, nil
}

func (l *LoggerImpl) Fatal(v ...any) {
	l.l.Fatal(v...)
}

func (l *LoggerImpl) Fatalf(format string, v ...any) {
	l.l.Fatalf(format, v...)
}
