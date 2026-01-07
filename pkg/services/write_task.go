package services

import (
	"os"
	"time"

	"encoding/csv"

	"github.com/fmo/timer-cli/pkg/models"
)

func WriteTask(file *os.File, task models.Task) error {
	start := task.Start.Format(time.RFC3339)
	var end string
	if task.End.IsZero() {
		end = ""
	} else {
		end = task.End.Format(time.RFC3339)
	}

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{start, end, string(task.Status)}); err != nil {
		return err
	}
	writer.Flush()

	return nil
}
