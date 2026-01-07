package services

import (
	"testing"

	"github.com/fmo/timer-cli/pkg/models"
)

func TestWriteTask(t *testing.T) {
	f, err := OpenFile("test.csv")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	task := models.NewTask()

	WriteTask(f, task)
}

func TestWriteTaskToEmptyFile(t *testing.T) {
	f, err := OpenFile("test.csv")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = f.Truncate(0)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	records, err := GetDataFromCSV(f)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(records) > 0 {
		t.Errorf("unexpected records")
	}

	task := models.NewTask()
	if err := WriteTask(f, task); err != nil {
		t.Errorf("unexpected error")
	}

	f, err = OpenFile("test.csv")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	records, err = GetDataFromCSV(f)
	if err != nil {
		t.Errorf("unexpected error")
	}

	tasks, err := GetTasks(records)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(tasks) < 1 {
		t.Error("should be one task here")
	}

	if !tasks[0].End.IsZero() {
		t.Errorf("end date should be empty for the new task")
	}
}
