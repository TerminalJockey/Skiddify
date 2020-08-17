package panels

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

//Loot returns the loot window back to main app
func Loot(window fyne.Window) fyne.CanvasObject {
	lootLabel := widget.NewLabel("loot tab, tbc")

	lootPage := fyne.NewContainerWithLayout(layout.NewCenterLayout(), lootLabel)
	return lootPage
}
