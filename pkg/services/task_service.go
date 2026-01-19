package services

import (
	"fmt"
	"log"
	"time"
)

type TaskService struct {
	Storer Storer
	Logger *log.Logger
	Tasks  *Tasks
}

func NewTaskService(s Storer, l *log.Logger) (*TaskService, error) {
	data, err := s.LoadData()
	if err != nil {
		return nil, err
	}

	tasks, err := NewTasks(data, l)
	if err != nil {
		return nil, err
	}

	return &TaskService{
		Storer: s,
		Logger: l,
		Tasks:  tasks,
	}, nil
}

func (ts *TaskService) Create() (*Task, error) {
	task := NewTask()
	task.Start()

	if err := ts.Tasks.AllowNewTask(); err != nil {
		return nil, err
	}

	if err := ts.Storer.Save(task); err != nil {
		return nil, err
	}

	return task, nil
}

func (ts *TaskService) AddManual(startTime, endTime time.Time) error {
	task := NewTask()
	task.StartTime = startTime
	task.Status = Done
	task.EndTime = endTime
	if err := ts.Storer.Save(task); err != nil {
		return err
	}
	return nil
}

func (ts *TaskService) Complete() error {
	currentTask, err := ts.GetCurrentTask()
	if err != nil {
		return fmt.Errorf("cant complete task due to: %w", err)
	}
	currentTask.Complete()

	return ts.Storer.Update(currentTask)
}

func (ts *TaskService) ResetData() error {
	return ts.Storer.ResetData()
}

func (ts *TaskService) GetCurrentTask() (*Task, error) {
	return ts.Tasks.GetCurrentTask()
}

func (ts *TaskService) TotalDuration() time.Duration {
	return ts.Tasks.TotalDuration()
}
