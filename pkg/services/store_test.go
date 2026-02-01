package services

import (
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

func TestStoreSave(t *testing.T) {
	store := newStore(t)

	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)

	task := &Task{
		StartTime: startTime,
		EndTime:   endTime,
		Status:    Started,
	}

	store.Save(task)

	records, err := store.LoadData()
	if err != nil {
		t.Error("unexpected error while loading data")
	}

	totalRecords := len(records)
	if totalRecords != 1 {
		t.Errorf("expected record number: %d, got: %d", 1, totalRecords)
	}
}
