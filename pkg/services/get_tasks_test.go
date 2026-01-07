package services

import (
	"testing"
	"time"
)

// Test creates time objects and converts them to string in RFC3339 format and pass it to the GetTasks
func TestGetTasks(t *testing.T) {
	start := time.Now()
	end := start.Add(30 * time.Second)

	// format times
	fstart := start.Format(time.RFC3339)
	fend := end.Format(time.RFC3339)

	data := [][]string{
		{"start", "end", "status"},
		{fstart, fend, "done"},
	}

	got, err := GetTasks(data)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if got[0].Start.Format(time.RFC3339) != fstart {
		t.Error("not matching start")
	}
	if got[0].End.Format(time.RFC3339) != fend {
		t.Error("not matching end")
	}
}

func TestGetTasksLessFields(t *testing.T) {
	data := [][]string{
		{"start", "end", "status"},
		{"", ""},
	}
	_, err := GetTasks(data)
	if err.Error() != "not matching column length" {
		t.Errorf("error expected here")
	}
}

func TestGetTasksWrongEnd(t *testing.T) {
	start := time.Now()
	end := start.Add(30 * time.Minute)

	fstart := start.Format(time.RFC3339)
	fend := end.Format("01-02")

	data := [][]string{
		{"start", "end", "status"},
		{fstart, fend, "done"},
	}
	_, err := GetTasks(data)
	if err == nil {
		t.Error("error expected here")
	}
}

func TestGetTasksWrongStatus(t *testing.T) {
	start := time.Now()
	end := start.Add(30 * time.Minute)
	fstart := start.Format(time.RFC3339)
	fend := end.Format(time.RFC3339)

	data := [][]string{
		{"start", "end", "status"},
		{fstart, fend, "hello"},
	}

	_, err := GetTasks(data)
	if err == nil {
		t.Error("error expected here")
	}
}
