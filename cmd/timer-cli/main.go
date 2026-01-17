package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fmo/timer-cli/pkg/services"
)

const (
	taskFile = "tasks.csv"
	layout   = "02-01-2006 15:04:05 MST"
)

func main() {
	// Logger
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("cant create log file")
	}

	logger := log.New(logFile, "logs: ", log.Lshortfile|log.LstdFlags)

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
		log.Fatal(err)
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
			log.Fatal(err)
		}
	case "complete":
		//currentTask, err := tasks.GetCurrentTask()
		//logger.Printf("Current task is: %v", currentTask)
		//currentTask.Store = taskStorer
		//if err != nil {
	//		logger.Fatalf("ther is no current task: %v", err)
	//	}
	//	if err := currentTask.Complete(); err != nil {
	//		logger.Fatalf("cant complete the task: %v", err)
	//	}
	case "add":
	case "show":
		currentTask, error := taskService.GetCurrentTask()
		if error != nil {
			log.Fatal(error)
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
