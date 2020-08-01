package scanpage

import (
	"fmt"
	"image/color"
	"net"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var line = canvas.NewLine(color.Black)
var space = layout.NewSpacer()
var cancelCounter = 0

// Scanner returns scanner
func Scanner(window fyne.Window) fyne.CanvasObject {

	/* port scanner */
	net := widget.NewLabel("Network Enumeration")

	/* build entrypoints for ip & range */
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("{Enter IP here}")
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("{Enter Port Range ie: 139-445}")

	/* build form and tie to entry points
	   introduces cancelCounter, when tapped sets
	   global variable to 1, include break statement
	   in long running functions to allow cancels
	   (see portScanner func) */

	ipForm := &widget.Form{
		Items: []*widget.FormItem{
			{"Enter target IP: ", ipEntry},
			{"Enter port range: ", portEntry}},
		OnSubmit: func() {
			ip := ipEntry.Text
			port := portEntry.Text
			fmt.Println("scan button tapped!")
			go PortScan(ip, port, window)
		},
		SubmitText: "Scan",
		OnCancel: func() {
			cancelCounter = 1
		},
		CancelText: "Cancel",
	}

	/* builds layout, first vertical box for forms then the whole page as horizontal box */
	netCol := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), net, ipForm, space, space)
	scanPage := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), line, netCol, space)
	return scanPage
}

//PortScan returns portscan
func PortScan(ip string, port string, window fyne.Window) {
	addr := strings.Split(ip, ".")
	ports := strings.Split(port, "-")

	/* check valid ip / port combo */
	if len(addr) < 4 {
		//errorPop(window, "IP address.")
		fmt.Println("invalid IP address.")
	}
	for i := range addr {
		num, _ := strconv.Atoi(addr[i])
		if num >= 255 || num < 0 {
			fmt.Println("invalid IP address.")
			//errorPop(window, "IP address.")
			break
		}
	}
	for j := range ports {
		hold, _ := strconv.Atoi(ports[j])
		if hold < 1 || hold > 65535 {
			fmt.Println("invalid port selection.")
			//errorPop(window, "port selection.")
			break
		}
	}

	/* handles single port provision */
	if strings.Contains(port, "-") == false {
		fmt.Println("no space found!")
		fmt.Println("scanning port: ", port)
		singleConn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), 500*time.Millisecond)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Port: ", port, "open")
			singleConn.Close()
		}

		/* handle port ranges */
	} else if strings.Contains(port, "-") == true {
		fmt.Println("Range identified, pulling...")
		p1, _ := strconv.Atoi(ports[0])
		p2, _ := strconv.Atoi(ports[1])
		if p1 < p2 {
			diff := p2 - p1
			for m := 0; m <= diff; m++ {
				//cancelCounter break statement
				if cancelCounter == 1 {
					cancelCounter = 0
					break
				}
				currPort := strconv.Itoa(p1 + m)
				conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, currPort), 500*time.Millisecond)
				if err != nil {
					fmt.Println("Port: ", currPort, "closed")
				} else {
					fmt.Println("Port:", currPort, "open")
					conn.Close()
				}
			}
		} else if p2 < p1 {
			diff := p1 - p2
			for m := 0; m <= diff; m++ {
				//cancelCounter break statement
				if cancelCounter == 1 {
					cancelCounter = 0
					break
				}
				currPort := strconv.Itoa(p2 + m)
				conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, currPort), 500*time.Millisecond)
				if err != nil {
					fmt.Println("Port: ", currPort, "closed")
				} else {
					fmt.Println("Port: ", currPort, "open")
					conn.Close()
				}
			}
		}
	}
}
