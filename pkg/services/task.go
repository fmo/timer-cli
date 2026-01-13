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
	Store  TaskStorer
	Start  time.Time
	End    time.Time
	Status Status
}

func NewTask(ts TaskStorer) *Task {
	return &Task{Store: ts}
}

func (t *Task) Create() error {
	t.Start = time.Now()
	t.Status = Started
	return t.Store.Save(*t)
}

func (t *Task) Complete() error {
	t.Status = Done
	t.End = time.Now()
	return t.Store.Update(t)
}
