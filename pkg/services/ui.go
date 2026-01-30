package services

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app          *tview.Application
	menu         *tview.Flex
	display      *tview.TextView
	displayPages *tview.Pages
	form         *tview.Form
	onSubmit     func(string, string)
	root         *tview.Flex
}

func NewUI() *UI {
	app := tview.NewApplication()

	menu := tview.NewFlex()
	menu.SetDirection(tview.FlexRow)
	menu.SetBorder(true)

	display := tview.NewTextView()
	display.SetBorder(true)
	display.SetText("Loading Tasks...")

	ui := &UI{
		app:     app,
		menu:    menu,
		display: display,
	}

	form := tview.NewForm()
	form.AddInputField("start time", "hh:mm:ss", 10, nil, nil)
	form.AddInputField("duration", "1h22m33s", 20, nil, nil)
	form.SetBorder(true)
	form.AddButton("Add Time", func() {
		st := form.GetFormItemByLabel("start time").(*tview.InputField).GetText()
		d := form.GetFormItemByLabel("duration").(*tview.InputField).GetText()

		if ui.onSubmit != nil {
			ui.onSubmit(st, d)
		}
		form.RemoveFormItem(0)
		form.RemoveFormItem(0)
		form.RemoveButton(0)
		form.AddTextView("Status", "Saved successfully", 40, 1, true, false)
	})

	ui.form = form

	ui.displayPages = tview.NewPages().
		AddPage("textBase", display, true, true).
		AddPage("form", form, true, false)

	content := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(ui.menu, 0, 1, true).
		AddItem(ui.displayPages, 0, 1, false)

	centered := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 1, false).
		AddItem(content, 80, 0, false).
		AddItem(nil, 0, 1, false)

	ui.root = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(centered, 23, 0, false).
		AddItem(nil, 0, 1, false)

	return ui
}

func (ui *UI) SubmitForm(fn func(string, string)) {
	ui.onSubmit = fn
}

func (ui *UI) SwitchToForm() {
	ui.displayPages.SwitchToPage("form")
}

func (ui *UI) SwitchToTextBase() {
	ui.displayPages.SwitchToPage("textBase")
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

func (ui *UI) Render() {
	ui.app.EnableMouse(true).SetRoot(ui.root, true).Run()
}

func (ui *UI) Stop() {
	ui.app.Stop()
}
