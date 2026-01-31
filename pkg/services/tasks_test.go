package services

import (
	"testing"
	"time"

	"github.com/fmo/timer-cli/pkg/logger"
)

func TestNewTasks(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(30 * time.Minute)

	task1 := []string{
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		"done",
	}

	task2 := []string{
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		"started",
	}

	logger, err := logger.New()
	if err != nil {
		t.Error("cant initiate the logger")
	}

	tasks, err := NewTasks([][]string{task1, task2}, logger)
	if err != nil {
		t.Error("cant create tasks")
	}

	if len(tasks.items) != 2 {
		t.Errorf("expected %d, got %d", 2, len(tasks.items))
	}
}
