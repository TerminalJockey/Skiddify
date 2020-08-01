package main

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

/* set global vars (mostly for formatting) */
var line = canvas.NewLine(color.Black)
var space = layout.NewSpacer()

func main() {
	/* initialize window and set size */
	app := app.NewWithID("Skiddify!")
	window := app.NewWindow("Skiddify!")
	window.Resize(fyne.NewSize(1200, 800))

	/* side menu ties to functions */
	sidebar := widget.NewTabContainer(
		widget.NewTabItem("Scanner", scanner(window)),
		widget.NewTabItem("Exploitation", exploitation()),
		widget.NewTabItem("Post", postEx()),
		widget.NewTabItem("Viewer", viewer()),
		widget.NewTabItem("Shell", shell()))

	sidebar.SetTabLocation(widget.TabLocationLeading)
	sidebar.SelectTabIndex(0)
	window.SetContent(sidebar)
	window.ShowAndRun()
}

/* scanner page */
func scanner(window fyne.Window) fyne.CanvasObject {

	/* port scanner */
	net := widget.NewLabel("Network Enumeration")

	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("{Enter IP here}")
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("{Enter Port Range ie: 139-445}")

	ipForm := &widget.Form{
		Items: []*widget.FormItem{
			{"Enter target IP: ", ipEntry},
			{"Enter port range: ", portEntry}},
		OnSubmit: func() {
			ip := ipEntry.Text
			port := portEntry.Text
			fmt.Println("scan button tapped!")
			portScan(ip, port, window)
		},
		SubmitText: "Scan",
	}

	netCol := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), net, ipForm, space, space)

	scanPage := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), line, netCol, space)

	return scanPage
}

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

func portScan(ip string, port string, window fyne.Window) {
	addr := strings.Split(ip, ".")
	ports := strings.Split(port, "-")

	//check valid ip / port combo
	if len(addr) < 4 {
		errorPop(window, "IP address.")
		fmt.Println("invalid IP address.")
	}
	for i := range addr {
		num, _ := strconv.Atoi(addr[i])
		if num >= 255 || num < 0 {
			fmt.Println("invalid IP address.")
			errorPop(window, "IP address.")
			break
		}
	}
	for j := range ports {
		hold, _ := strconv.Atoi(ports[j])
		if hold < 1 || hold > 65535 {
			fmt.Println("invalid port selection.")
			errorPop(window, "port selection.")
			break
		}
	}

}

func errorPop(window fyne.Window, message string) {
	errMsg := "Error! Check " + message
	test := canvas.NewText(errMsg, color.Black)
	content := fyne.NewContainerWithLayout(layout.NewCenterLayout(), test)
	canvas := window.Canvas()
	widget.ShowPopUpAtPosition(content, canvas, fyne.NewPos(600, 400))
}
