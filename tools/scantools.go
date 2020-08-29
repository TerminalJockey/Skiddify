package tools

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//PortScan scans ports lol
func PortScan(IP string, PortRange string, threads int) {
	fmt.Println("scan started:", IP, PortRange)

	ports := formatPorts(PortRange)

	go pScan(ports, IP, PortRange, threads)
}

func formatPorts(PortRange string) (portArray []string) {

	clearSpace := strings.ReplaceAll(PortRange, " ", "")
	if strings.Contains(clearSpace, "-") == true || strings.Contains(clearSpace, ",") == true {
		ports := strings.Split(PortRange, ",")
		for i := 0; i < len(ports); i++ {
			if strings.Contains(ports[i], "-") == true {
				takSeparated := strings.Split(ports[i], "-")
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
			}
		}
	} else {
		portArray = append(portArray, PortRange)
	}
	return portArray
}

//takes PortRange to show the requested scan in string format at output
func pScan(ports []string, IP string, PortRange string, threads int) {

	if IP == "" {
		IP = "localhost"
	}

	wg := sync.WaitGroup{}

	var openPorts []string

	var sem = make(chan int, threads)
	//sem channel acts as throttle, good spot for improvement
	if threads == 0 {
		sem = nil
	}

	//res recieves results from goroutine
	res := make(chan string)

	//AHAHAHA VICTORY
	for _, port := range ports {
		//add throttle counter
		if sem != nil {
			sem <- 1
		}
		wg.Add(1)
		defer wg.Done()
		go scanPort(IP, port, res, &wg, sem)

	}

	for range ports {
		openPorts = append(openPorts, <-res)
	}

	//open|create results file
	resFile, err := os.OpenFile("PortscanResults.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer resFile.Close()

	if err != nil {
		log.Println(err)
	}

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	//host format header:
	header := fmt.Sprintf("Finished scan of host: %s ports: %s \n\n\n", IP, PortRange)
	if _, err := resFile.Write([]byte(header)); err != nil {
		resFile.Close()
		log.Fatal(err)
	}
	for x := range openPorts {
		if openPorts[x] != "" {
			result := fmt.Sprintf("Port: %s is open\n", openPorts[x])
			if _, err := resFile.Write([]byte(result)); err != nil {
				resFile.Close()
				log.Fatal(err)
			}
		}
	}

	if _, err := resFile.Write([]byte("\n\n\n")); err != nil {
		resFile.Close()
		log.Fatal(err)
	}

	fmt.Println("Done with pScan")

}

func scanPort(IP string, port string, res chan string, wg *sync.WaitGroup, sem chan int) {
	//get open ports
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(IP, port), time.Millisecond*300)
	//release counter
	if sem != nil {
		<-sem
	}
	if err != nil {
		log.Println(err)
		res <- ""
	} else {
		conn.Close()
		res <- port
	}
}

//ClearResults deletes a given file
func ClearResults(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Println(err)
	}
}

//DirScan scans dirs...
func DirScan() {
	fmt.Println("dirscan")
}

//BruteForce gon' do some bruteforcin'
func BruteForce() {
	fmt.Println("bruteforce")
}
