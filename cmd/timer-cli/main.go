package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fmo/timer-cli/pkg/logger"
	"github.com/fmo/timer-cli/pkg/services"
)

type tickMsg string

type taskStateMsg struct {
	isRunning bool
	total     string
}

var docStyle = lipgloss.NewStyle().MarginTop(20).MarginLeft(56)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list   list.Model
	width  int
	height int

	taskService *services.TaskService

	timerCtx    context.Context
	cancelTimer context.CancelFunc
	elapsed     string
	total       string
	isRunning   bool
}

func (m model) Init() tea.Cmd {
	return initTaskState(m.taskService)
}

func (m *model) ensureTimer() context.Context {
	if m.timerCtx == nil {
		ctx, cancel := context.WithCancel(context.Background())
		m.timerCtx = ctx
		m.cancelTimer = cancel
	}
	return m.timerCtx
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if it, ok := m.list.SelectedItem().(item); ok {
				if it.title == "Start" {
					var task *services.Task
					var err error

					if m.isRunning {
						task, err = m.taskService.GetCurrentTask()
						if err != nil {
							m.elapsed = err.Error()
						}
					} else {
						task, err = m.taskService.Create()
						if err != nil {
							m.elapsed = err.Error()
							return m, nil
						}
						m.isRunning = true
					}

					ctx := m.ensureTimer()

					return m, countTime(ctx, task)
				}
				if it.title == "Complete" {
					if m.cancelTimer != nil {
						m.cancelTimer()
					}

					m.timerCtx = nil
					m.cancelTimer = nil

					if err := m.taskService.Complete(); err != nil {
						m.elapsed = err.Error()
						return m, nil
					}
					m.isRunning = false
					m.elapsed = "0s"
					return m, nil
				}
				if it.title == "Show" {
					currentTask, err := m.taskService.GetCurrentTask()
					if err != nil {
						m.elapsed = err.Error()
						return m, nil
					}
					ctx := m.ensureTimer()
					return m, countTime(ctx, currentTask)
				}
				if it.title == "Total" {
					m.total = m.taskService.TotalDuration()
					return m, nil
				}
				if it.title == "Reset" {
					if m.cancelTimer != nil {
						m.cancelTimer()
					}

					m.timerCtx = nil
					m.cancelTimer = nil

					if err := m.taskService.ResetData(); err != nil {
						m.elapsed = err.Error()
					}

					return m, nil
				}
				if it.title == "Close" {
					return m, tea.Quit
				}
			}
		case "ctrl+c":
			if m.cancelTimer != nil {
				m.cancelTimer()
			}
			return m, tea.Quit
		}
	case tickMsg:
		m.elapsed = string(msg)
		currentTask, _ := m.taskService.GetCurrentTask()
		return m, countTime(m.timerCtx, currentTask)
	case taskStateMsg:
		m.isRunning = msg.isRunning
		m.total = msg.total
		if m.isRunning {
			currentTask, _ := m.taskService.GetCurrentTask()
			ctx := m.ensureTimer()
			return m, countTime(ctx, currentTask)
		}

		return m, nil
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

var (
	leftStyle  = lipgloss.NewStyle().Width(40).Padding(0, 1)
	rightStyle = lipgloss.NewStyle().Width(70).Padding(0, 1)
)

func rightView(m model) string {
	if m.list.SelectedItem() == nil {
		return "Select a menu item"
	}

	i := m.list.SelectedItem().(item)

	switch i.title {
	case "Start":
		if m.isRunning {
			return "Elapsed:\n\n" + m.elapsed
		}
		return "Start a new task"
	case "Show":
		if m.isRunning {
			return "Elapsed:\n\n" + m.elapsed
		}
		return "There is no started task yet"
	case "Total":
		return "Total:\n\n" + m.total
	case "Complete":
		if m.isRunning {
			return "Complete the task"
		}
		return "Task completed"
	case "Reset":
		return "Reset the tasks"
	case "Close":
		return "Close the app"
	default:
		return ""
	}
}

func (m model) View() string {
	left := leftStyle.Render(m.list.View())
	right := rightStyle.Render(rightView(m))

	ui := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	return docStyle.Render(ui)
}

func initTaskState(ts *services.TaskService) tea.Cmd {
	return func() tea.Msg {
		td := ts.TotalDuration()
		_, err := ts.GetCurrentTask()
		if err != nil {
			return taskStateMsg{isRunning: false, total: td}
		}

		return taskStateMsg{isRunning: true, total: td}
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

	items := []list.Item{
		item{title: "Start", desc: "Start the task"},
		item{title: "Show", desc: "Show running task"},
		item{title: "Complete", desc: "Complete the task"},
		item{title: "Total", desc: "Total todays tasks"},
		item{title: "Reset", desc: "Reset csv"},
		item{title: "Close", desc: "Closing the app"},
	}

	m := model{
		list:        list.New(items, list.NewDefaultDelegate(), 0, 0),
		taskService: taskService,
	}
	m.list.Title = "My Fave Things"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func countTime(ctx context.Context, task *services.Task) tea.Cmd {
	return func() tea.Msg {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second):
			elapsed := time.Since(task.StartTime).
				Truncate(time.Second).
				String()
			return tickMsg(elapsed)
		}
	}
}
