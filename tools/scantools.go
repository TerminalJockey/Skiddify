package tools

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type pScanObject struct {
	IP        string
	Portrange string
	Ports     []string
	Threads   int
}

type pScanResult struct {
	IP     string
	Port   string
	State  string
	Banner string
}

type dScanObject struct {
	Host       string
	extensions string
	extArray   []string
	wordlist   string
	directory  string
	url        string
	threads    int
	waitTime   int
	IsTLS      bool
}

type dScanResult struct {
	Host       string
	statusCode string
	directory  string
}

func InitPortScan(IP string, Portrange string, Threads int) {

	if IP == "" {
		IP = "localhost"
	}

	var scanObject pScanObject
	var scanResult pScanResult
	var resultChannel = make(chan pScanResult)

	scanObject.IP = strings.TrimSpace(IP)
	scanObject.Portrange = Portrange
	scanObject.Ports = getPorts(strings.TrimSpace(Portrange))
	scanObject.Threads = Threads
	var sem = make(chan int)
	var wg sync.WaitGroup

	sem = semSetup(Threads)

	for _, port := range scanObject.Ports {
		if sem != nil {
			sem <- 1
		}
		wg.Add(1)
		defer wg.Done()
		go portScan(scanObject, port, scanResult, sem, &wg, resultChannel)
	}

	fmt.Println("passed goroutine")

	var scanArray []pScanResult
	for range scanObject.Ports {
		scanArray = append(scanArray, <-resultChannel)
	}

	if len(scanArray) > 0 {
		writeResult(scanArray)
	}
}

func portScan(scanObject pScanObject, port string, scanResult pScanResult, sem chan int, wg *sync.WaitGroup, resultChannel chan pScanResult) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(scanObject.IP, port), time.Millisecond*300)
	if sem != nil {
		<-sem
	}
	if err != nil {
		scanResult.IP = scanObject.IP
		scanResult.Port = port
		scanResult.State = "closed"
		resultChannel <- scanResult
	} else {
		conn.Close()
		scanResult.IP = scanObject.IP
		scanResult.Port = port
		scanResult.State = "open"
		resultChannel <- scanResult
	}

}

func getPorts(Portrange string) (portArray []string) {
	if strings.Contains(Portrange, "-") == true || strings.Contains(Portrange, ",") == true {
		ports := strings.Split(Portrange, ",")
		for i := 0; i < len(ports); i++ {
			if strings.Contains(ports[i], "-") == true {
				takSeparated := strings.Split(ports[i], "-")
				val0, err := strconv.Atoi(takSeparated[0])
				handleError(err)
				val1, err := strconv.Atoi(takSeparated[1])
				handleError(err)
				if val0 > val1 {
					for x := 0; x < ((val0 - val1) + 1); x++ {
						portArray = append(portArray, strconv.Itoa(val1+x))
					}
				} else if val1 > val0 {
					for x := 0; x < ((val1 - val0) + 1); x++ {
						portArray = append(portArray, strconv.Itoa(val0+x))
					}
				}
			}
		}
	} else {
		portArray = append(portArray, Portrange)
	}
	return portArray
}

func InitDirScan(IP string, extensions string, threads int, wordlist string) {
	var scanObject dScanObject
	var scanResult dScanResult
	var resultChannel = make(chan dScanResult)
	var counter int
	var sem = make(chan int)
	var wg sync.WaitGroup

	scanObject.Host = strings.TrimSpace(IP)
	scanObject.threads = threads
	scanObject.wordlist = wordlist
	scanObject.IsTLS = checkTLS(scanObject.Host)
	scanObject.extensions = strings.TrimSpace(extensions)

	if extensions != "" {
		scanObject.extArray = strings.Split(scanObject.extensions, ",")
	} else {
		scanObject.extArray = nil
	}

	sem = semSetup(threads)

	file, err := os.Open(scanObject.wordlist)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	if scanObject.extArray != nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			counter++
			for _, ext := range scanObject.extArray {
				if sem != nil {
					sem <- 1
				}
				if scanObject.IsTLS == true {
					scanObject.url = "https://" + scanObject.Host + "/" + scanner.Text() + "." + ext
					wg.Add(1)
					defer wg.Done()
					go scanDir(scanObject, scanResult, sem, &wg, resultChannel)
				} else {
					scanObject.url = "http://" + scanObject.Host + "/" + scanner.Text() + "." + ext
					wg.Add(1)
					defer wg.Done()
					go scanDir(scanObject, scanResult, sem, &wg, resultChannel)
				}
			}
		}
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			counter++
			if sem != nil {
				sem <- 1
			}
			if scanObject.IsTLS == true {
				scanObject.url = "https://" + scanObject.Host + "/" + scanner.Text()
				wg.Add(1)
				defer wg.Done()
				go scanDir(scanObject, scanResult, sem, &wg, resultChannel)
			} else {
				scanObject.url = "http://" + scanObject.Host + "/" + scanner.Text()
				wg.Add(1)
				defer wg.Done()
				go scanDir(scanObject, scanResult, sem, &wg, resultChannel)
			}
		}
	}
	var scanArray []dScanResult
	if counter == 0 {
		counter = 1
	}

	for x := 0; x < (counter * len(scanObject.extArray)); x++ {
		scanArray = append(scanArray, <-resultChannel)
	}
	if len(scanArray) > 0 {
		writeResult(scanArray)
	}

}

func scanDir(scanObject dScanObject, scanResult dScanResult, sem chan int, wg *sync.WaitGroup, resultChannel chan dScanResult) {
	client := http.Client{Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	get, err1 := http.NewRequest("GET", scanObject.url, nil)
	get.Close = true
	resp, err2 := client.Do(get)
	if sem != nil {
		<-sem
	}
	handleError(err1)
	handleError(err2)
	if err1 == nil && err2 == nil {
		scanResult.directory = scanObject.url
		scanResult.statusCode = resp.Status
		scanResult.Host = scanObject.Host
		resultChannel <- scanResult
	} else {
		scanResult.directory = scanObject.directory
		scanResult.statusCode = "err"
		scanResult.Host = "err"
		resultChannel <- scanResult
	}
}

func checkTLS(IP string) bool {
	_, err := http.Get(IP)
	if err != nil {
		if strings.HasSuffix(err.Error(), "x509: certificate signed by unknown authority") == true {
			return true
		}
	}
	return false
}

func semSetup(Threads int) (sem chan int) {
	if Threads == 0 {
		sem = make(chan int)
		sem = nil
		return sem
	}
	sem = make(chan int, Threads)
	return sem

}

func handleError(err error) {
	if err != nil {
		log.Println(err)
	}
}

//ClearResults deletes a given file
func ClearResults(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Println(err)
	}
}

func writeResult(x interface{}) {

	resultFile, err := os.OpenFile("ScanResults.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer resultFile.Close()
	handleError(err)

	switch y := x.(type) {
	case []pScanResult:
		header := fmt.Sprintf("Finished Portscan \n\n\n")
		if _, err := resultFile.Write([]byte(header)); err != nil {
			resultFile.Close()
			log.Fatal(err)
		}
		for z := range y {
			if y[z].Port != "" {
				write := fmt.Sprintf("%s:%s is %s\n", y[z].IP, y[z].Port, y[z].State)
				if _, err := resultFile.Write([]byte(write)); err != nil {
					resultFile.Close()
					log.Fatalln(err)
				}
			}
		}
	case []dScanResult:
		header := fmt.Sprintf("Finished DirScan of host: %s \n\n\n", y[1].Host)
		if _, err := resultFile.Write([]byte(header)); err != nil {
			resultFile.Close()
			log.Fatal(err)
		}
		for v := range y {
			fmt.Println(y[v].directory)
			if y[v].statusCode != "404 Not Found" && y[v].statusCode != "err" {
				write := fmt.Sprintf("%s %s \n", y[v].directory, y[v].statusCode)
				if _, err := resultFile.Write([]byte(write)); err != nil {
					resultFile.Close()
					log.Fatal(err)
				}
			}
		}
	}
	if _, err := resultFile.Write([]byte("\n\n")); err != nil {
		resultFile.Close()
		log.Fatal(err)
	}
}
