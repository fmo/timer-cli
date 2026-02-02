package services

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fmo/timer-cli/pkg/logger"
)

func newStore(t *testing.T) Storer {
	logger, err := logger.New()
	if err != nil {
		t.Error("cant initiate logger")
	}

	codec, err := NewCSVCodec(logger)
	if err != nil {
		t.Error("cant initiate the csv codec")
	}

	codec.ResetData()

	return NewStore(codec)
}

func taskData() []*Task {
	var tasks []*Task

	for range 5 {
		duration := rand.Intn(200)
		startTime := time.Now()
		endTime := startTime.Add(time.Duration(duration) * time.Minute)

		task := &Task{
			StartTime: startTime,
			EndTime:   endTime,
			Status:    Started,
		}

		tasks = append(tasks, task)
	}

	return tasks
}

func TestStoreSave(t *testing.T) {
	store := newStore(t)

	taskData := taskData()

	store.Save(taskData[0])

	records, err := store.LoadData()
	if err != nil {
		t.Error("unexpected error while loading data")
	}

	totalRecords := len(records)
	if totalRecords != 1 {
		t.Errorf("expected record number: %d, got: %d", 1, totalRecords)
	}
}

func TestStoreUpdate(t *testing.T) {
	store := newStore(t)
	taskData := taskData()

	if err := store.Save(taskData[0]); err != nil {
		t.Errorf("unexpected save error: %v", err)
	}

	taskData[0].Status = Done

	if err := store.Update(taskData[0]); err != nil {
		t.Errorf("unexpected update error: %v", err)
	}

	storedData, err := store.LoadData()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	totalDoneTask := 0
	for _, d := range storedData {
		if d[2] == "done" {
			totalDoneTask++
		}
	}

	if totalDoneTask != 1 {
		t.Errorf("expected: %d, got: %d", 1, totalDoneTask)
	}
}
