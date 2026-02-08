package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fmo/timer-cli/pkg/logger"
	"github.com/fmo/timer-cli/pkg/services"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
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

	err        string
	success    string
	manualMode bool
	cursorMode cursor.Mode
	focusIndex int
	inputs     []textinput.Model
}

// Init things
func (m model) Init() tea.Cmd {
	return tea.Batch(
		initTaskState(m.taskService),
		textinput.Blink,
	)
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

func (m *model) initManualInputs() {
	m.inputs = make([]textinput.Model, 2)

	for i := range m.inputs {
		t := textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Start Time (11:00:00)"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Duration (1h30m10s)"
			t.CharLimit = 64
		}
		m.inputs[i] = t
	}
}

func (m *model) ensureTimer() context.Context {
	if m.timerCtx == nil {
		ctx, cancel := context.WithCancel(context.Background())
		m.timerCtx = ctx
		m.cancelTimer = cancel
	}
	return m.timerCtx
}

func (m model) updateManual(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.manualMode = false
			return m, nil
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > cursor.CursorHide {
				m.cursorMode = cursor.CursorBlink
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].Cursor.SetMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				startTime, err := stringToTime(m.inputs[0].Value())
				if err != nil {
					m.err = err.Error()
					return m, nil
				}
				duration, err := time.ParseDuration(m.inputs[1].Value())
				if err != nil {
					m.err = err.Error()
					return m, nil
				}
				endTime := startTime.Add(duration)
				m.taskService.AddManual(startTime, endTime)
				m.success = "task added successfully"

				m.total = m.taskService.TotalDuration()
				m.manualMode = false
				return m, nil
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
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

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.manualMode {
		return m.updateManual(msg)
	}

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
					m.total = m.taskService.TotalDuration()
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
				if it.title == "Manual" {
					m.manualMode = true
					m.focusIndex = 0

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

// View section
func (m model) View() string {
	left := leftStyle.Render(m.list.View())
	right := rightStyle.Render(rightView(m))

	ui := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	return docStyle.Render(ui)
}

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("9")).
	Bold(true)

func (m model) manualView() string {
	var b strings.Builder

	if m.err != "" {
		b.WriteString(errorStyle.Render(m.err))
		b.WriteString("\n\n")
	}

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	b.WriteString("\n\n[esc] back")

	return b.String()
}

func rightView(m model) string {
	if m.manualMode {
		return m.manualView()
	}

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
	case "Manual":
		if m.success != "" {
			return m.success
		}
		return "Add manual task"
	case "Close":
		return "Close the app"
	default:
		return ""
	}
}

// Main section
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
		item{title: "Manual", desc: "Add manual task"},
		item{title: "Reset", desc: "Reset csv"},
		item{title: "Close", desc: "Closing the app"},
	}

	m := model{
		list:        list.New(items, list.NewDefaultDelegate(), 0, 0),
		taskService: taskService,
	}
	m.list.Title = "Cli Task Manager"
	m.initManualInputs()

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
