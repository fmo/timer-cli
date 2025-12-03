package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	taskFile = "tasks.csv"
	layout   = "02-01-2006 15:04:05 MST"
)

type Task struct {
	Start  *time.Time
	End    *time.Time
	Status string // started, done
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("need at least one argument"))
	}

	f, err := getCsv()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	tasks, err := getTasks(f)
	if err != nil {
		log.Fatal(err)
	}

	currentTask := getCurrentTask(tasks)

	switch os.Args[1] {
	case "start":
		s := time.Now()
		if err := write(f, Task{Start: &s, End: nil, Status: "started"}); err != nil {
			log.Fatal(err)
		}
	case "total":
        fmt.Printf("Total time: %v\n", total(tasks))        
    case "reset":
		if err := resetData(f); err != nil {
			log.Fatal(err)
		}
	case "end":
		currentTask.Status = "done"
		now := time.Now()
		currentTask.End = &now
		updateTask(f, tasks, currentTask)
	case "add":
        if len(os.Args) < 4 {
            log.Fatal(errors.New("not enough arguments"))
        }
		addStart := os.Args[2]
		addDuration := os.Args[3]
		now := time.Now()
		location := now.Location()

		startArr := strings.Split(addStart, ":")

		if len(startArr) < 3 {
            log.Fatal(errors.New("Not enough arguments"))
        }

        hour, _ := strconv.Atoi(startArr[0])
		minute, _ := strconv.Atoi(startArr[1])
		sec, _ := strconv.Atoi(startArr[2])

		startDate := time.Date(now.Year(), time.Month(12), now.Day(), hour, minute, sec, int(00), location)

        d, _ := strconv.Atoi(addDuration)
        
        endDate := startDate.Add(time.Minute * time.Duration(d))

        record := Task{Start: &startDate, End: &endDate, Status: "done"}
		addManual(f, record)
	default:
		if currentTask == nil {
			fmt.Printf("No active task running\n")
			os.Exit(0)
		}
		countTime(currentTask)

	}
}

func total(tasks []Task) time.Duration {
    var d time.Duration
    for _, task := range tasks {
        start := *task.Start
        end := *task.End
        diff := end.Sub(start)
        d += diff
    }
    return d
} 

func resetData(f *os.File) error {
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	if err := writeHeader(f); err != nil {
		log.Fatal(err)
	}

	return nil
}

func countTime(task *Task) error {
	if task == nil {
		return errors.New("empty task")
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var d time.Duration
	for range ticker.C {
		now := time.Now()
		d = now.Sub(*task.Start)
		d = d.Truncate(time.Second)
		fmt.Print("\033[H\033[2J")
		fmt.Println(d.String())
	}

	return nil
}

func getCsv() (*os.File, error) {
	var f *os.File
	var err error

	f, err = os.OpenFile(taskFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func write(f *os.File, task Task) error {
	w := csv.NewWriter(f)
	s := task.Start.Format(layout)
	record := []string{s, "", task.Status}
	w.Write(record)
	w.Flush()
	return nil
}

func writeHeader(f *os.File) error {
	w := csv.NewWriter(f)

	record := []string{"start", "end", "status"}

	if err := w.Write(record); err != nil {
		return err
	}

	w.Flush()
	return nil
}

func updateTask(f *os.File, tasks []Task, task *Task) error {
	updatedTasks := []Task{}
	for _, t := range tasks {
		if t.Start == task.Start {
			updatedTasks = append(updatedTasks, *task)
			continue
		}
		updatedTasks = append(updatedTasks, t)
	}

	resetData(f)

	w := csv.NewWriter(f)

	for _, t := range updatedTasks {
		s := t.Start.Format(layout)
		e := t.End.Format(layout)
		record := []string{s, e, t.Status}

		if err := w.Write(record); err != nil {
			return err
		}
		w.Flush()
	}

	return nil
}

func getTasks(f *os.File) ([]Task, error) {
	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for i := 1; i < len(data); i++ {
		var parsedStart, parsedEnd time.Time
		var err error

		if data[i][0] != "" {
			parsedStart, err = time.Parse(layout, data[i][0])
			if err != nil {
				return nil, err
			}
		}

		if data[i][1] != "" {
			parsedEnd, err = time.Parse(layout, data[i][1])
			if err != nil {
				return nil, err
			}
		}

		tasks = append(tasks, Task{Start: &parsedStart, End: &parsedEnd, Status: data[i][2]})
	}

	return tasks, nil
}

func getCurrentTask(tasks []Task) *Task {
	for _, task := range tasks {
		if task.Status == "started" {
			return &task
		}
	}

	return nil
}

func addManual(f *os.File, t Task) {
	w := csv.NewWriter(f)
	start := t.Start.Format(layout)
	end := t.End.Format(layout)
	record := []string{start, end, "done"}
	w.Write(record)
	w.Flush()
}
