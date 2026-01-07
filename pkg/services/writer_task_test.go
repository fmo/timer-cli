package services

import (
	"testing"

	"github.com/fmo/timer-cli/pkg/models"
)

func TestWriteTask(t *testing.T) {
	f, err := OpenFile("test1.csv")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	task := models.NewTask()

	WriteTask(f, task)
}
