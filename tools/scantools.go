package tools

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

//PortScan scans ports lol
func PortScan(IP string, PortRange string) []string {
	fmt.Println("scan started:", IP, PortRange)
	clean := strings.ReplaceAll(PortRange, " ", "")

	var ports []string

	if strings.Contains(clean, "-") == true || strings.Contains(clean, ",") == true {
		ports = getRange(PortRange)
	} else {
		ports = append(ports, PortRange)
	}

	for i := range ports {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(IP, ports[i]), 500*time.Millisecond)
		if err != nil {
			continue
		} else {
			fmt.Printf("Host: %s Port: %s is open\n", IP, ports[i])
			conn.Close()
		}
	}
	return ports
}

func getRange(PortRange string) []string {

	var portArray []string
	separated := strings.Split(PortRange, ",")

	for i := 0; i < len(separated); i++ {
		if strings.Contains(separated[i], "-") == true {

			takSeparated := strings.Split(separated[i], "-")
			val0, err := strconv.Atoi(takSeparated[0])
			if err != nil {
				log.Println(err)
			}
			val1, err := strconv.Atoi(takSeparated[1])
			if err != nil {
				log.Println(err)
			}

			if val0 > val1 {
				diff := val0 - val1
				for x := 0; x < diff+1; x++ {
					appendVal := strconv.Itoa(val1 + x)
					portArray = append(portArray, appendVal)
				}
			} else if val1 > val0 {
				diff := val1 - val0
				for x := 0; x < diff+1; x++ {
					appendVal := strconv.Itoa(val0 + x)
					portArray = append(portArray, appendVal)
				}
			}

		} else {
			portArray = append(portArray, separated[i])
		}
	}
	return portArray
}

//DirScan scans dirs...
func DirScan() {
	fmt.Println("dirscan")
}

//BruteForce gon' do some bruteforcin'
func BruteForce() {
	fmt.Println("bruteforce")
}
