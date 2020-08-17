package panels

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

//Advanced returns the advanced window back to main app
func Advanced(window fyne.Window) fyne.CanvasObject {
	advancedLabel := widget.NewLabel("Advanced tab, tbc")

	advancedPage := fyne.NewContainerWithLayout(layout.NewCenterLayout(), advancedLabel)
	return advancedPage
}
