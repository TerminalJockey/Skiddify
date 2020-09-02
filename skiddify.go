package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/TerminalJockey/Skiddify/0.3/panels"
)

func main() {
	app := app.NewWithID("Skiddify")
	window := app.NewWindow("Skiddify")
	window.Resize(fyne.NewSize(1200, 800))

	sidebar := widget.NewTabContainer(
		widget.NewTabItem("Scan", panels.Scan(window)),
		widget.NewTabItem("Attack", panels.Attack(window)),
		widget.NewTabItem("Loot", panels.Loot(window)),
		widget.NewTabItem("Control", panels.Control(window)),
		widget.NewTabItem("Advanced", panels.Advanced(window)))

	sidebar.SetTabLocation(widget.TabLocationLeading)
	sidebar.SelectTabIndex(0)
	window.SetContent(sidebar)
	window.ShowAndRun()
}
