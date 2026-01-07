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

	"github.com/fmo/timer-cli/pkg/models"
	"github.com/fmo/timer-cli/pkg/services"
)

const (
	taskFile = "tasks.csv"
	layout   = "02-01-2006 15:04:05 MST"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need at least one argument")
	}

	file, err := services.OpenFile(taskFile)
	if err != nil {
		log.Fatal(err)
	}

	data, err := services.GetDataFromCSV(file)
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := services.GetTasks(data)
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "start":
		task := models.NewTask()
		if err := write(f, task); err != nil {
			log.Fatal(err)
		}
		countTime(task)
	case "total":
		fmt.Printf("Total time: %v\n", total(tasks))
	case "reset":
		if err := resetData(f); err != nil {
			log.Fatal(err)
		}
	case "end":
		currentTask, err := services.GetCurrentTask(tasks)
		if err != nil {
			log.Fatal(err)
		}
		currentTask.Status = models.Done
		now := time.Now()
		currentTask.End = now
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
			log.Fatal(errors.New("not enough arguments"))
		}

		hour, _ := strconv.Atoi(startArr[0])
		minute, _ := strconv.Atoi(startArr[1])
		sec, _ := strconv.Atoi(startArr[2])

		startDate := time.Date(now.Year(), time.Month(12), now.Day(), hour, minute, sec, int(00), location)

		d, _ := strconv.Atoi(addDuration)

		endDate := startDate.Add(time.Minute * time.Duration(d))

		record := models.Task{Start: &startDate, End: &endDate, Status: models.Done}
		addManual(f, record)
	case "show":
		if currentTask == nil {
			fmt.Printf("No active task running\n")
			os.Exit(0)
		}
		countTime(currentTask)
	default:
		fmt.Println("Usage: ")
		fmt.Println("  timer-cli <command>")
		fmt.Println("")
		fmt.Println("Commands:")
		fmt.Println("  timer-cli start -- starts the task")
		fmt.Println("  timer-cli end -- ends the task")
		fmt.Println("  timer-cli total -- total time during day")
		fmt.Println("  timer-cli reset -- reset the whole file and adds the header to csv")
		fmt.Println("  timer-cli add -- adds manual time")
		fmt.Println("  timer-cli show -- shows the current active task's running time")
	}
}

func total(tasks []models.Task) time.Duration {
	var d time.Duration
	day := time.Now().Format("02")
	for _, task := range tasks {
		start := task.Start
		taskDay := start.Format("02")
		if taskDay != day {
			continue
		}
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

func countTime(task *models.Task) error {
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

func writeHeader(f *os.File) error {
	w := csv.NewWriter(f)

	record := []string{"start", "end", "status"}

	if err := w.Write(record); err != nil {
		return err
	}

	w.Flush()
	return nil
}

func updateTask(f *os.File, tasks []models.Task, task *models.Task) error {
	updatedTasks := []models.Task{}
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
		record := []string{s, e, string(t.Status)}

		if err := w.Write(record); err != nil {
			return err
		}
		w.Flush()
	}

	return nil
}

func addManual(f *os.File, t models.Task) {
	w := csv.NewWriter(f)
	start := t.Start.Format(layout)
	end := t.End.Format(layout)
	record := []string{start, end, string(models.Done)}
	w.Write(record)
	w.Flush()
}
