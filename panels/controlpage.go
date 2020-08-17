package panels

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

//Control returns the control window back to main app
func Control(window fyne.Window) fyne.CanvasObject {
	controlLabel := widget.NewLabel("Control tab, tbc")

	page := fyne.NewContainerWithLayout(layout.NewCenterLayout(), controlLabel)
	return page
}
