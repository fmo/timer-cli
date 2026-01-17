package services

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

type Task struct {
	StartTime time.Time
	EndTime   time.Time
	Status    Status
}

func NewTask() *Task {
	return &Task{}
}

func (t *Task) Start() {
	t.StartTime = time.Now()
	t.Status = Started
}

func (t *Task) HasStarted() bool {
	return t.Status == Started
}

func (t *Task) Complete() {
	t.Status = Done
	t.EndTime = time.Now()
}

func (t *Task) IsSameTask(startTime string) bool {
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return false
	}
	return t.StartTime.Equal(start)
}
