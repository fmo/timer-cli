package services

import (
	"fmt"
	"strconv"
	"strings"
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

// AddManual Data came as;
// StarTime 11:00:00
// Duration 1h20m30s
func (ts *TaskService) AddManual(st, d string) error {
	startTime, err := stringToTime(st)
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(d)
	if err != nil {
		return err
	}

	endTime := startTime.Add(duration)

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

func (ts *TaskService) TotalDuration() string {
	totalDuration := ts.Tasks.TotalDuration()
	totalDuration = totalDuration.Truncate(1 * time.Second)
	return totalDuration.String()
}

func stringToTime(s string) (time.Time, error) {
	timeArr := strings.Split(s, ":")
	if len(timeArr) < 3 {
		return time.Time{}, fmt.Errorf("need starting time format hh:mm::ss")
	}

	hh, mm, ss := timeArr[0], timeArr[1], timeArr[2]
	hhInt, err := strconv.Atoi(hh)
	if err != nil {
		return time.Time{}, fmt.Errorf("cant get the hour")
	}
	mmInt, err := strconv.Atoi(mm)
	if err != nil {
		return time.Time{}, fmt.Errorf("cant get the minute")
	}
	ssInt, err := strconv.Atoi(ss)
	if err != nil {
		return time.Time{}, fmt.Errorf("cant get the second")
	}

	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), hhInt, mmInt, ssInt, 0, now.Location())

	return t, nil
}
