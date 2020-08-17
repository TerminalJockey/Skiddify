package panels

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"
)

//Scan returns the scan window back to main app
func Scan(window fyne.Window) fyne.CanvasObject {

	line := canvas.NewLine(color.Black)
	line1 := canvas.NewLine(color.Black)

	//make labels for funcitonality
	pScannerLabel := widget.NewLabel("Port Scanner")
	dScannerLabel := widget.NewLabel("Directory Scanner")
	bruteForceLabel := widget.NewLabel("Brute Forcer")
	resultsLabel := widget.NewLabel("Results")

	//portscanbox
	pScanIPEntry := widget.NewEntry()
	pScanIPEntry.SetPlaceHolder("{Target IP}")
	pScanPortEntry := widget.NewEntry()
	pScanPortEntry.SetPlaceHolder("{Port Range}")

	//build portscan entry forms
	pScanForm := &widget.Form{
		Items: []*widget.FormItem{
			{"Target IP: ", pScanIPEntry},
			{"Port Range: ", pScanPortEntry}},
		OnSubmit: func() {
			pTargIP := pScanIPEntry.Text
			pTargPorts := pScanPortEntry.Text
			fmt.Println("scan started:", pTargIP, pTargPorts)
		},
		SubmitText: "Scan",
		OnCancel: func() {
			fmt.Println("portscan cancelled!")
		},
		CancelText: "Cancel",
	}

	//assemble portscan vertical box
	vPScanBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), pScannerLabel, pScanForm)

	//assemble portscan horizontal box
	hPScanBox := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), vPScanBox)

	//dirscanbox
	dScanIPEntry := widget.NewEntry()
	dScanIPEntry.SetPlaceHolder("Enter URL/IP")
	dScanExtEntry := widget.NewEntry()
	dScanExtEntry.SetPlaceHolder("Enter comma separated extensions")

	//wordlist selection

	wlButton := widget.NewButton("Select Wordlist", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			try := strings.Split(fmt.Sprintf("%s", reader), " ")
			path := strings.TrimSuffix(try[len(try)-1], "}")
			fmt.Println(path)
		}, window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		fd.Show()
	})

	//build dirscan entry forms
	dScanForm := &widget.Form{
		Items: []*widget.FormItem{
			{"Target IP/URL: ", dScanIPEntry},
			{"Extensions: ", dScanExtEntry}},
		OnSubmit: func() {
			dScanTargIP := dScanIPEntry.Text
			dScanExt := dScanExtEntry.Text
			fmt.Println("dir scan started: ", dScanTargIP, dScanExt)
		},
		SubmitText: "DirScan",
		OnCancel: func() {
			fmt.Println("dirscan cancelled!")
		},
		CancelText: "Cancel",
	}

	//assemble dirscan vertical box
	vDScanBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), dScannerLabel, wlButton, dScanForm)
	//assemble dirscan horizontal box
	hDScanBox := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), line, vDScanBox)

	//assemble top horizontal area (port scanner and dir scanner)
	topH := fyne.NewContainerWithLayout(layout.NewGridLayout(2), hPScanBox, hDScanBox)

	topV := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), topH, line1)
	page := fyne.NewContainerWithLayout(layout.NewGridLayout(1), topV, bruteForceLabel, layout.NewSpacer(), resultsLabel)

	return page

}
