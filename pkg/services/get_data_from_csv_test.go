package services

import (
	"testing"
)

func TestGetDataFromCSV(t *testing.T) {
	f, err := OpenFile("test.csv")
	if err != nil {
		t.Error("unexpected error", err)
	}
	_, err = GetDataFromCSV(f)
	if err != nil {
		t.Error("unexpected error", err)
	}
}
