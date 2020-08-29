package panels

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/widget"

	"github.com/TerminalJockey/Skiddify/0.2/tools"
)

//Scan returns the scan window back to main app
func Scan(window fyne.Window) fyne.CanvasObject {

	tDiv := canvas.NewLine(color.Black)
	bDiv := canvas.NewLine(color.Black)

	//make labels for funcitonality
	pScannerLabel := widget.NewLabel("Port Scanner")
	dScannerLabel := widget.NewLabel("Directory Scanner")
	bruteForceLabel := widget.NewLabel("Brute Forcer")

	//portscanbox
	pScanIPEntry := widget.NewEntry()
	pScanIPEntry.SetPlaceHolder("Target IP")
	pScanPortEntry := widget.NewEntry()
	pScanPortEntry.SetPlaceHolder("Port Range")
	pScanThreads := widget.NewEntry()
	pScanThreads.SetPlaceHolder("Scan Threads")

	//build portscan entry forms
	pScanForm := &widget.Form{
		Items: []*widget.FormItem{
			{"Target IP: ", pScanIPEntry},
			{"Port Range: ", pScanPortEntry},
			{"Threads: ", pScanThreads}},
		OnSubmit: func() {
			pTargIP := pScanIPEntry.Text
			pTargPorts := pScanPortEntry.Text
			pTargThreads, _ := strconv.Atoi(pScanThreads.Text)
			go tools.PortScan(pTargIP, pTargPorts, pTargThreads)
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
		wfd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			wPath := getPath(reader)
			fmt.Println(wPath)
		}, window)
		wfd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		wfd.Show()
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
	hDScanBox := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), vDScanBox)

	//bruteforcer
	bForceIPEntry := widget.NewEntry()
	bForceIPEntry.SetPlaceHolder("Target IP")
	bForcePortEntry := widget.NewEntry()
	bForcePortEntry.SetPlaceHolder("Port")

	//ip|port submission form

	bForceForm := &widget.Form{
		Items: []*widget.FormItem{
			{"Target IP: ", bForceIPEntry},
			{"Target Port: ", bForcePortEntry}},
		OnSubmit: func() {
			bForceIP := bForceIPEntry.Text
			bForcePort := bForceIPEntry.Text
			fmt.Printf("bruteforce started with options: ip: %s port: %s \n", bForceIP, bForcePort)
		},
		SubmitText: "BruteForce",
		OnCancel: func() {
			fmt.Println("bruteforce cancelled")
		},
		CancelText: "Cancel",
	}

	//user list
	uList := widget.NewButton("Select Userlist", func() {
		ufd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			uPath := getPath(reader)
			fmt.Println("upath", uPath)
		}, window)
		ufd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		ufd.Show()
	})

	//pass list
	pList := widget.NewButton("Select PassList", func() {
		pfd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			pPath := getPath(reader)
			fmt.Println(pPath)
		}, window)
		pfd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		pfd.Show()
	})

	//setup lefthand side of bruteforce panel
	LLPanelbForceBox := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), bForceForm)
	LRPanelbForceBox := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(120, 40)), layout.NewSpacer(), uList, pList)

	bForceBox := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), LLPanelbForceBox, layout.NewSpacer(), LRPanelbForceBox, layout.NewSpacer())

	//bForce Service selection panel

	bForceSelectionLabel := widget.NewLabel("Select Service")

	//create service buttons
	smbButton := widget.NewButton("SMB", func() {
		fmt.Println("smb tapped!")
	})

	ftpButton := widget.NewButton("FTP", func() {
		fmt.Println("ftp tapped")
	})

	sshButton := widget.NewButton("SSH", func() {
		fmt.Println("ssh tapped")
	})

	sftpButton := widget.NewButton("SFTP", func() {
		fmt.Println("sftp tapped")
	})

	ldapButton := widget.NewButton("LDAP", func() {
		fmt.Println("ldap tapped")
	})

	smtpButton := widget.NewButton("SMTP", func() {
		fmt.Println("smtp tapped")
	})

	mysqlButton := widget.NewButton("MYSQL", func() {
		fmt.Println("mysql tapped")
	})

	mssqlButton := widget.NewButton("MSSQL", func() {
		fmt.Println("mssql tapped")
	})

	rdpButton := widget.NewButton("RDP", func() {
		fmt.Println("rdp tapped")
	})

	//assemble buttons into 3x3 grid, with spacers to keep the grid from the side of the window
	bForceSelectionGrid := fyne.NewContainerWithLayout(layout.NewGridLayout(4), smbButton, ftpButton, rdpButton, layout.NewSpacer(),
		sshButton, sftpButton, ldapButton, layout.NewSpacer(), smtpButton, mysqlButton, mssqlButton, layout.NewSpacer())

	//assemble bruteforce selection panel
	bForceSelectionPanel := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), bForceSelectionLabel, bForceSelectionGrid)

	//begin results box
	resultsContent, err := ioutil.ReadFile("PortscanResults.txt")
	if err != nil {
		log.Println(err)
	}

	//without passing text to a widget.NewMultiLineEntry, scroll box becomes horribly laggy
	resultsText := widget.NewMultiLineEntry()
	if string(resultsContent) == "" {
		resultsText.SetText("No results yet...")
	} else {
		resultsText.SetText(string(resultsContent))
		resultsText.Wrapping = fyne.TextWrapBreak
	}

	//allow results tab to refresh, might be nice to have an auto-update but on demand seems fine
	refreshButton := widget.NewButton("Refresh Results", func() {
		resultsContent, err = ioutil.ReadFile("PortscanResults.txt")
		if err != nil {
			log.Println(err)
		}
		resultsText.SetText(string(resultsContent))
		resultsText.Wrapping = fyne.TextWrapBreak
		resultsText.Refresh()
	})

	clearButton := widget.NewButton("Clear Results", func() {
		tools.ClearResults("PortscanResults.txt")
	})

	resultsScroll := widget.NewVScrollContainer(resultsText)
	resSize := fyne.NewSize(1200, 220)
	resultsScroll.SetMinSize(resSize)
	resultsScroll.Refresh()

	accord := widget.NewAccordionItem("Results", resultsScroll)
	resultsRender := widget.NewAccordionContainer(accord)

	buttonBar := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), refreshButton, clearButton)

	resultsBar := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), resultsRender, buttonBar)

	resultsBox := fyne.NewContainerWithLayout(layout.NewAdaptiveGridLayout(1), resultsBar)

	//assemble top horizontal area (port scanner and dir scanner)
	topH := fyne.NewContainerWithLayout(layout.NewAdaptiveGridLayout(2), hPScanBox, hDScanBox)

	//assemble mid horizontal area (bruteforce entry form, userlist/passlist select, and service selection grid)
	midH := fyne.NewContainerWithLayout(layout.NewAdaptiveGridLayout(2), bForceBox, bForceSelectionPanel)

	//assemble bottom horizontal area (results box)
	botH := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), resultsBox)

	//organize panels
	topV := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), topH, layout.NewSpacer(), tDiv)
	midV := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), bruteForceLabel, layout.NewSpacer(), midH, layout.NewSpacer(), bDiv)

	//create page
	page := fyne.NewContainerWithLayout(layout.NewGridLayout(1), topV, midV, botH)

	return page

}

//func getPath retrieves the filepath of a target from the fyne.URIReadCloser, allowing us to pull words from wordlists
func getPath(reader fyne.URIReadCloser) string {
	raw := strings.Split(fmt.Sprintf("%s", reader), " ")
	path := strings.TrimSuffix(raw[len(raw)-1], "}")
	return path
}
