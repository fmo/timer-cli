package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fmo/timer-cli/pkg/logger"
	"github.com/fmo/timer-cli/pkg/services"
)

const taskFile = "tasks.csv"

func main() {
	// Logger
	logger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		logger.Fatal("need at least one argument")
	}

	// File for CSV
	file, err := os.OpenFile(taskFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Fatalf("cant open the file: %v", err)
	}

	// CSV Codec
	persister, err := services.NewCSVCodec(file, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Storer
	storer := services.NewStore(persister)

	// Task Service
	taskService, err := services.NewTaskService(storer, logger)
	if err != nil {
		logger.Fatal(err)
	}

	switch os.Args[1] {
	case "start":
		task, err := taskService.Create()
		if err != nil {
			logger.Fatal(err)
		}
		countTime(task)
	case "total":
		fmt.Printf("Total time: %v\n", taskService.TotalDuration())
	case "reset":
		if err := taskService.ResetData(); err != nil {
			logger.Fatal(err)
		}
	case "complete":
		if err := taskService.Complete(); err != nil {
			logger.Fatal(err)
		}
	case "add":
		if len(os.Args) < 4 {
			logger.Fatal("need start time and duration for manual adding")
		}

		startTimeInString := os.Args[2]

		startTime, err := stringTimeToTime(startTimeInString)
		if err != nil {
			log.Fatal(err)
		}

		addition := os.Args[3]
		additionInt, err := strconv.Atoi(addition)
		if err != nil {
			log.Fatal("need addition time")
		}

		endTime := startTime.Add(time.Duration(additionInt) * time.Minute)

		taskService.AddManual(startTime, endTime)
	case "show":
		currentTask, error := taskService.GetCurrentTask()
		if error != nil {
			logger.Fatal(error)
		}
		countTime(currentTask)
	case "help":
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

func stringTimeToTime(s string) (time.Time, error) {
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

func countTime(task *services.Task) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var d time.Duration
	for range ticker.C {
		now := time.Now()
		d = now.Sub(task.StartTime)
		d = d.Truncate(time.Second)
		fmt.Print("\033[H\033[2J")
		fmt.Println(d.String())
	}

	return nil
}
