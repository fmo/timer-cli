package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/fmo/timer-cli/pkg/logger"
)

type Tasks struct {
	logger logger.Logger
	items  []Task
}

func NewTasks(data [][]string, logger logger.Logger) (*Tasks, error) {
	tasks := &Tasks{logger: logger}

	items, err := tasks.getAll(data)
	if err != nil {
		return nil, err
	}

	tasks.items = items

	return tasks, nil
}

func (t *Tasks) AddTask(task Task) {
	t.items = append(t.items, task)
}

func (t *Tasks) UpdateTask(task Task) {
	for i, v := range t.items {
		if v.StartTime.Equal(task.StartTime) {
			t.items[i].Status = Done
			t.items[i].EndTime = task.EndTime
		}
	}
}

func (t *Tasks) RemoveAll() {
	t.items = []Task{}
}

func (t *Tasks) TotalDuration() time.Duration {
	var total time.Duration
	for _, task := range t.items {
		if task.IsTodaysTask() && task.HasDone() {
			total += task.Duration()
		}
	}

	return total
}

func (t *Tasks) AllowNewTask() error {
	for _, task := range t.items {
		if task.HasStarted() {
			return errors.New("a task already running")
		}
	}
	return nil
}

func (t *Tasks) getAll(tasksArr [][]string) ([]Task, error) {
	var tasks []Task

	for _, taskArr := range tasksArr {
		if len(taskArr) != 3 {
			return nil, errors.New("not matching column length")
		}
		task := Task{}

		start, err := time.Parse(time.RFC3339, taskArr[0])
		if err != nil {
			return nil, fmt.Errorf("parsing start time not possible: %w", err)
		}
		task.StartTime = start

		if taskArr[1] != "" {
			end, err := time.Parse(time.RFC3339, taskArr[1])
			if err != nil {
				return nil, fmt.Errorf("parsing end time not possible: %w", err)
			}
			task.EndTime = end
		}
		status := Status(taskArr[2])
		if err := status.IsValid(); err != nil {
			return nil, fmt.Errorf("cant assign status: %w", err)
		}
		task.Status = status
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (t *Tasks) GetCurrentTask() (*Task, error) {
	for _, task := range t.items {
		if task.HasStarted() {
			return &task, nil
		}
	}
	return nil, errors.New("there is no started task")
}
