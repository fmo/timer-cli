package services

import (
	"testing"
	"time"

	"github.com/fmo/timer-cli/pkg/models"
)

func TestRemoveCSVHeader(t *testing.T) {
	data := [][]string{
		{
			time.Now().Format(time.RFC3339),
			time.Now().Add(10 * time.Minute).Format(time.RFC3339),
			string(models.Done),
		},
	}

	dataRemovedHeader := RemoveCSVHeader(data)

	if dataRemovedHeader[0][2] != string(models.Done) {
		t.Error("unexpected response")
	}
}
