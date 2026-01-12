package main

import (
	"errors"
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
	if len(os.Args) < 2 {
		log.Fatal("need at least one argument")
	}

	file, err := os.OpenFile(taskFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("cant open the file: %v", err)
	}

	csvCodec, err := services.NewCSVCodec(file)
	if err != nil {
		if errors.Is(err, services.ErrNoDataInCSV) {
			log.Println("no data in csv")
		} else {
			log.Fatalf("cant get the codec: %v", err)
		}
	}

	data, err := csvCodec.Load()
	if err != nil {
		if errors.Is(err, services.ErrNoDataInCSV) {
			log.Println("csv is empty")
		} else {
			log.Fatalf("something is wrong while loading csv: %v", err)
		}
	}

	tasks, err := services.NewTasks(data)
	if err != nil {
		log.Fatalf("cant get the tasks object: %v", err)
	}

	taskStorer := services.NewStore(csvCodec)

	switch os.Args[1] {
	case "start":
		task := services.NewTask(taskStorer)
		if err := task.Create(); err != nil {
			log.Fatal(err)
		}
		countTime(*task)
	case "total":
		fmt.Printf("Total time: %v\n", tasks.TotalDuration())
	case "reset":
		if err := csvCodec.ResetData(); err != nil {
			log.Fatal(err)
		}
	case "complete":
		currentTask, err := tasks.GetCurrentTask()
		if err != nil {
			log.Fatalf("there is no current task: %v", err)
		}
		if err := currentTask.Complete(); err != nil {
			log.Fatalf("cant complete the task: %v", err)
		}
	case "add":
	case "show":
		currentTask, error := tasks.GetCurrentTask()
		if error != nil {
			log.Fatal("there is no current task")
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

func countTime(task services.Task) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var d time.Duration
	for range ticker.C {
		now := time.Now()
		d = now.Sub(task.Start)
		d = d.Truncate(time.Second)
		fmt.Print("\033[H\033[2J")
		fmt.Println(d.String())
	}

	return nil
}
