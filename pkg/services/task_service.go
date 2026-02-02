package services

import (
	"fmt"
	"time"

	"github.com/fmo/timer-cli/pkg/logger"
)

type TaskService struct {
	Storer Storer
	Logger logger.Logger
	Tasks  *Tasks
}

func NewTaskService(s Storer, l logger.Logger) (*TaskService, error) {
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

	ts.Tasks.AddTask(task)

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

	ts.Tasks.AddTask(task)

	return nil
}

func (ts *TaskService) Complete() error {
	currentTask, err := ts.GetCurrentTask()
	if err != nil {
		return fmt.Errorf("cant complete task due to: %v", err)
	}

	currentTask.Complete()

	ts.Tasks.UpdateTask(currentTask)

	return ts.Storer.Update(currentTask)
}

func (ts *TaskService) ResetData() error {
	if err := ts.Storer.ResetData(); err != nil {
		return err
	}

	ts.Tasks.RemoveAll()

	return nil
}

func (ts *TaskService) GetCurrentTask() (*Task, error) {
	return ts.Tasks.GetCurrentTask()
}

func (ts *TaskService) TotalDuration() time.Duration {
	return ts.Tasks.TotalDuration()
}
