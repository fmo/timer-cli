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
		log.Fatal(err, "err")
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
		_, err := taskService.Create()
		if err != nil {
			logger.Fatal(err)
		}
		//countTime(task)
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

		startTime, err := stringToTime(startTimeInString)
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
		_, error := taskService.GetCurrentTask()
		if error != nil {
			logger.Fatal(error)
		}
		//countTime(currentTask)
	case "help":
		ui := services.NewUI()
		f := func() {
			task, err := taskService.Create()
			if err != nil {
				ui.SetDisplayText(err.Error())
				return
			}
			go countTime(task, func(text string) {
				ui.SetDynamicDisplayText(text)
			})
		}
		ui.AddMenuItem("start", "start the task", f)
		ui.DrawLayout()
	}
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

func countTime(task *services.Task, update func(string)) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		elapsed := time.Since(task.StartTime).
			Truncate(1 * time.Second).
			String()
		update(elapsed)
	}
}
