package main

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	scanpage "github.com/TerminalJockey/Skiddify/panels"
)

/* set global vars (mostly for formatting) */
var line = canvas.NewLine(color.Black)
var space = layout.NewSpacer()
var cancelCounter = 0

func main() {
	/* initialize window and set size */
	app := app.NewWithID("Skiddify!")
	window := app.NewWindow("Skiddify!")
	window.Resize(fyne.NewSize(1200, 800))

	/* side menu ties to functions */
	sidebar := widget.NewTabContainer(
		widget.NewTabItem("Scanner", scanpage.Scanner(window)),
		widget.NewTabItem("Exploitation", exploitation()),
		widget.NewTabItem("Post", postEx()),
		widget.NewTabItem("Viewer", viewer()),
		widget.NewTabItem("Shell", shell()))

	sidebar.SetTabLocation(widget.TabLocationLeading)
	sidebar.SelectTabIndex(0)
	window.SetContent(sidebar)
	window.ShowAndRun()
}

//WILL BE MIGRATING TO PACKAGES

/* exploitation tab */
func exploitation() fyne.CanvasObject {
	tester := widget.NewLabel("testing extab")
	return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, tester, nil, nil), tester)
}

/* post exploitation tab (get dat loot) */
func postEx() fyne.CanvasObject {
	tester := widget.NewLabel("testing post")
	return fyne.NewContainerWithLayout(layout.NewVBoxLayout(), tester)
}

/* viewer tab, display compromised targets & info */
func viewer() fyne.CanvasObject {
	tester := widget.NewLabel("testing viewer")
	return fyne.NewContainerWithLayout(layout.NewCenterLayout(), tester)
}

/* shell access to host|revshells */
func shell() fyne.CanvasObject {
	tester := widget.NewLabel("testing shell")
	return fyne.NewContainerWithLayout(layout.NewCenterLayout(), tester)
}

/* handle error messages, create popup with recommendation to remediate */
func errorPop(window fyne.Window, message string) {
	errMsg := "Error! Check " + message
	test := canvas.NewText(errMsg, color.Black)
	content := fyne.NewContainerWithLayout(layout.NewCenterLayout(), test)
	canvas := window.Canvas()
	widget.ShowPopUpAtPosition(content, canvas, fyne.NewPos(600, 400))
}
