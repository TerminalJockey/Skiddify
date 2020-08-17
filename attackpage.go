package panels

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

//Attack returns the attack window to main app
func Attack(window fyne.Window) fyne.CanvasObject {
	attackLabel := widget.NewLabel("Attack tab, tbc")

	page := fyne.NewContainerWithLayout(layout.NewCenterLayout(), attackLabel)
	return page
}
