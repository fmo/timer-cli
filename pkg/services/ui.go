package services

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app     *tview.Application
	menu    *tview.Flex
	display *tview.TextView
}

func NewUI() *UI {
	app := tview.NewApplication()

	menu := tview.NewFlex()
	menu.SetDirection(tview.FlexRow)
	menu.SetBorder(true)

	display := tview.NewTextView()
	display.SetBorder(true)
	display.SetText("Loading Tasks...")

	return &UI{app, menu, display}
}

func (ui *UI) AddMenuItem(label, desc string, fn func()) {
	btn := tview.NewButton(label)
	btn.SetSelectedFunc(fn)

	btn.SetStyle(
		tcell.StyleDefault.
			Background(tcell.NewHexColor(0x1e3a8a)).
			Foreground(tcell.NewHexColor(0xe5e7eb)),
	)

	btn.SetActivatedStyle(
		tcell.StyleDefault.
			Background(tcell.NewHexColor(0x3b82f6)).
			Foreground(tcell.ColorWhite),
	)

	ui.menu.AddItem(btn, 3, 0, true)
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

	ui.app.SetFocus(ui.menu)
	ui.app.EnableMouse(true).SetRoot(root, true).Run()
}

func (ui *UI) Stop() {
	ui.app.Stop()
}
