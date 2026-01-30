package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fmo/timer-cli/pkg/logger"
	"github.com/fmo/timer-cli/pkg/services"
)

var cancelTimer context.CancelFunc

func stopTimer() {
	if cancelTimer != nil {
		cancelTimer()
		cancelTimer = nil
	}
}

func main() {
	// Logger
	logger, err := logger.New()
	if err != nil {
		log.Fatal(err, "err")
	}

	// CSV Codec
	persister, err := services.NewCSVCodec(logger)
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

	// Initiate UI
	ui := services.NewUI()

	startFn := func() {
		stopTimer()
		task, err := taskService.Create()
		if err != nil {
			ui.SetDisplayText(err.Error())
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancelTimer = cancel
		go countTime(ctx, task, func(text string) {
			ui.SetDynamicDisplayText(text)
		})
	}
	completeFn := func() {
		stopTimer()
		if err := taskService.Complete(); err != nil {
			ui.SetDisplayText(err.Error())
			return
		}
		ui.SetDisplayText("task completed")
	}
	showFn := func() {
		stopTimer()
		currentTask, err := taskService.GetCurrentTask()
		if err != nil {
			ui.SetDisplayText(err.Error())
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancelTimer = cancel
		go countTime(ctx, currentTask, func(text string) {
			ui.SetDynamicDisplayText(text)
		})
	}
	totalFn := func() {
		ui.SwitchToTextBase()
		stopTimer()
		totalDuration := taskService.TotalDuration()
		totalDuration = totalDuration.Truncate(1 * time.Second)
		ui.SetDisplayText(totalDuration.String())
	}
	resetFn := func() {
		stopTimer()
		if err := taskService.ResetData(); err != nil {
			ui.SetDisplayText(err.Error())
			return
		}
		ui.SetDisplayText("reset done")
	}
	manualFn := func() {
		ui.SubmitForm(func(st, d string) {
			startTime, err := stringToTime(st)
			if err != nil {
				logger.Fatal(err)
			}
			duration, err := time.ParseDuration(d)
			if err != nil {
				logger.Fatal(err)
			}
			endTime := startTime.Add(duration)
			if err := taskService.AddManual(startTime, endTime); err != nil {
				logger.Fatal(err)
			}
		})
		ui.SwitchToForm()
	}
	closeFn := func() {
		ui.Stop()
	}

	ui.AddMenuItem("start", "start the task", startFn)
	ui.AddMenuItem("complete", "complete the task", completeFn)
	ui.AddMenuItem("show", "show running task", showFn)
	ui.AddMenuItem("total", "show total duration", totalFn)
	ui.AddMenuItem("reset", "reset the data", resetFn)
	ui.AddMenuItem("manual", "add manual task", manualFn)
	ui.AddMenuItem("close", "close the timer", closeFn)

	// Default show the running task
	showFn()

	ui.Render()
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

func countTime(ctx context.Context, task *services.Task, update func(string)) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			elapsed := time.Since(task.StartTime).
				Truncate(1 * time.Second).
				String()
			update(elapsed)
		}
	}

}
