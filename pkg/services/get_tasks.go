// Package services
package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/fmo/timer-cli/pkg/models"
)

// GetTasks gets the string time in RFC3339 format and parses it to the object
// Also stripts the header part out
func GetTasks(data [][]string) ([]models.Task, error) {
	var tasks []models.Task

	tasksArr := RemoveCSVHeader(data)

	for _, taskArr := range tasksArr {
		if len(taskArr) != 3 {
			return nil, errors.New("not matching column length")
		}
		task := models.Task{}
		start, err := time.Parse(time.RFC3339, taskArr[0])
		if err != nil {
			return nil, fmt.Errorf("parsing start time not possible: %w", err)
		}
		task.Start = start

		if taskArr[1] != "" {

			end, err := time.Parse(time.RFC3339, taskArr[1])
			if err != nil {
				return nil, fmt.Errorf("parsing end time not possible: %w", err)
			}
			task.End = end
		}
		status := models.Status(taskArr[2])
		if err := status.IsValid(); err != nil {
			return nil, fmt.Errorf("cant assign status: %w", err)
		}
		task.Status = status
		tasks = append(tasks, task)
	}
	return tasks, nil
}
