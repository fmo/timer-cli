package services

import (
	"github.com/rivo/tview"
)

type UI struct {
	app       *tview.Application
	menuItems *tview.List
}

func NewUI() *UI {
	app := tview.NewApplication()
	menuItems := tview.NewList()
	menuItems.SetBorder(true)
	return &UI{app, menuItems}
}

func (ui *UI) AddMenuItem(label, desc string, fn func()) {
	ui.menuItems.AddItem(label, desc, 0, fn)
}

func (ui *UI) DrawLayout() {
	right := tview.NewTextView()
	right.SetBorder(true)
	right.SetText("right")

	content := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(ui.menuItems, 0, 1, true).
		AddItem(right, 0, 1, false)

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
