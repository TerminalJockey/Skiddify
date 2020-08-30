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

//PortScan scans ports lol
func PortScan(IP string, PortRange string, threads int) {
	fmt.Println("scan started:", IP, PortRange)

	ports := formatPorts(PortRange)

	IP = strings.TrimSpace(IP)
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

/* first we will verify the connection to the url. if good connection,
read lines from wordlist file, append words to url, append extension to url
and loop through combinations. return output to slice, write slice to results
as with portscanner. Output should be response code, and url. lets try using
structs this time.
*/

type urlDir struct {
	url      string
	response string
}

//DirScan scans dirs...
func DirScan(IP string, extensions string, threads int, wordlist string) {
	targ := strings.TrimSpace(IP)
	extensions = strings.TrimSpace(extensions)
	url := "http://" + targ
	//test connection to url
	testConn, err := http.Get(url)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(testConn.Status)
		go dScan(IP, extensions, threads, wordlist)
	}
}

func dScan(IP string, extensions string, threads int, wordlist string) {
	file, err := os.Open(wordlist)
	if err != nil {
		log.Println(err)
	}
	extArray := strings.Split(extensions, ",")

	defer file.Close()
	scanner := bufio.NewScanner(file)

	dResults := make(chan urlDir)
	dWG := sync.WaitGroup{}
	dSem := make(chan int, threads)
	if threads == 0 {
		dSem = nil
	}

	wListLen := 0
	for scanner.Scan() {
		wListLen++
		for _, ext := range extArray {
			if dSem != nil {
				dSem <- 1
			}
			url := "http://" + IP + "/" + scanner.Text() + "." + ext
			fmt.Println(url)
			dWG.Add(1)
			defer dWG.Done()
			go grabStatus(url, dResults, dSem)

		}
	}

	var codes []urlDir
	for x := 0; x < (wListLen * len(extArray)); x++ {
		codes = append(codes, <-dResults)
	}

	resFile, err := os.OpenFile("PortscanResults.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer resFile.Close()

	if err != nil {
		log.Println(err)
	}

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	//host format header:
	header := fmt.Sprintf("Finished DirScan of host: %s extensions: %s \n\n\n", IP, extensions)
	if _, err := resFile.Write([]byte(header)); err != nil {
		resFile.Close()
		log.Fatal(err)
	}
	for x := range codes {
		if codes[x].response != "404 Not Found" {
			result := fmt.Sprintf("%s %s \n", codes[x].response, codes[x].url)
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
}

func grabStatus(url string, dResults chan urlDir, dSem chan int) {
	var grabResult urlDir
	grab, err := http.Get(url)
	if dSem != nil {
		<-dSem
	}
	if err != nil {
		log.Println(err)
	}
	grabResult.response = grab.Status
	grabResult.url = url
	dResults <- grabResult

}

//BruteForce gon' do some bruteforcin'
func BruteForce() {
	fmt.Println("bruteforce")
}
