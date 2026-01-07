// Package models for whole project wide models
package models

import (
	"errors"
	"time"
)

var ErrWrongStatus = errors.New("wrong status")

type Status string

func (s Status) IsValid() error {
	if s != Started && s != Done {
		return ErrWrongStatus
	}
	return nil
}

const (
	Started Status = "started"
	Done    Status = "done"
)

// Task object keeps the task related data
type Task struct {
	Start  time.Time
	End    time.Time
	Status Status
}

func NewTask() Task {
	now := time.Now()
	return Task{Start: now, Status: Started}
}
