package services

import (
	"errors"

	"github.com/fmo/timer-cli/pkg/models"
)

func GetCurrentTask(tasks []models.Task) (models.Task, error) {
	for _, task := range tasks {
		if task.Status == models.Started {
			return task, nil
		}
	}
	return models.Task{}, errors.New("there is no started task")
}
