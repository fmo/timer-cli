package services

import (
	"github.com/rivo/tview"
)

type UI struct {
	app     *tview.Application
	menu    *tview.List
	display *tview.TextView
}

func NewUI() *UI {
	app := tview.NewApplication()

	menu := tview.NewList()
	menu.SetBorder(true)

	display := tview.NewTextView()
	display.SetBorder(true)
	display.SetText("right")

	return &UI{app, menu, display}
}

func (ui *UI) AddMenuItem(label, desc string, fn func()) {
	ui.menu.AddItem(label, desc, 0, fn)
}

func (ui *UI) SetDisplayText(text string) {
	ui.display.SetText(text)
}

func (ui *UI) SetDynamicDisplayText(text string) {
	ui.app.QueueUpdateDraw(func() {
		ui.display.SetText(text)
	})
}

func (ui *UI) DrawLayout() {
	content := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(ui.menu, 0, 1, true).
		AddItem(ui.display, 0, 1, false)

	centered := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(content, 80, 0, false).
		AddItem(nil, 0, 1, false)

	root := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(centered, 20, 0, false).
		AddItem(nil, 0, 1, false)

	ui.app.EnableMouse(true).SetRoot(root, true).Run()
}
