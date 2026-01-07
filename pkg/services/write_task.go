package services

import (
	"os"
	"time"

	"encoding/csv"

	"github.com/fmo/timer-cli/pkg/models"
)

func WriteTask(file *os.File, task models.Task) error {
	start := task.Start.Format(time.RFC3339)
	end := task.End.Format(time.RFC3339)
	writer := csv.NewWriter(file)
	writer.Write([]string{start, end, string(task.Status)})
	writer.Flush()

	return nil
}

