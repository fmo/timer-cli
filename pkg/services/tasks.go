package services

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type Tasks struct {
	Logger *log.Logger
	Items  []Task
}

func NewTasks(data [][]string, logger *log.Logger) (*Tasks, error) {
	tasks := &Tasks{
		Logger: logger,
	}
	items, err := tasks.GetAll(data)
	if err != nil {
		return nil, err
	}
	tasks.Items = items
	return tasks, nil
}

func (t *Tasks) TotalDuration() time.Duration {
	today := time.Now()
	day := today.Format("02")
	var total time.Duration
	for _, task := range t.Items {
		taskDay := task.Start.Format("02")
		if taskDay == day && task.Status == Done {
			total += task.End.Sub(task.Start)
		}

	}

	return total
}

func (t *Tasks) GetAll(tasksArr [][]string) ([]Task, error) {
	var tasks []Task

	for _, taskArr := range tasksArr {
		if len(taskArr) != 3 {
			return nil, errors.New("not matching column length")
		}
		t.Logger.Println(taskArr)
		task := Task{}

		start, err := time.Parse(time.RFC3339, taskArr[0])
		if err != nil {
			return nil, fmt.Errorf("parsing start time not possible: %w", err)
		}
		task.Start = start

		if taskArr[1] != "" {
			t.Logger.Printf("task row's first column is: %v", taskArr[1])
			end, err := time.Parse(time.RFC3339, taskArr[1])
			if err != nil {
				return nil, fmt.Errorf("parsing end time not possible: %w", err)
			}
			task.End = end
		}
		status := Status(taskArr[2])
		if err := status.IsValid(); err != nil {
			return nil, fmt.Errorf("cant assign status: %w", err)
		}
		task.Status = status
		t.Logger.Printf("task is %v", task)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (t *Tasks) GetCurrentTask() (*Task, error) {
	for _, task := range t.Items {
		if task.Status == Started {
			return &task, nil
		}
	}
	return &Task{}, errors.New("there is no started task")
}
