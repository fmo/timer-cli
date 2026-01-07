package models

import "testing"

func TestNewTask(t *testing.T) {
	task := NewTask()
	if !task.End.IsZero() {
		t.Error("Should be zero for the ending time for the task")
	}
}
