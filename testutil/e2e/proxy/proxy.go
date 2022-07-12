package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var mockFolder string = "testutil/e2e/proxy/mockMaps/"

var responsesChanged bool = false
var realCount int = 0
var cacheCount int = 0
var fakeCount int = 0

var fakeResponse bool = false

var saveJsonEvery int = 10 // in seconds
var epochTime int = 3      // in seconds
var epochCount int = 0     // starting epoch count
var proxies []proxyProcess = []proxyProcess{}

type proxyProcess struct {
	id        string
	port      string
	host      string
	mockfile  string
	mock      *mockMap
	handler   func(http.ResponseWriter, *http.Request)
	malicious bool
	cache     bool
	strict    bool
	noSave    bool
}

func getDomain(s string) (domain string) {
	parts := strings.Split(s, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return s
}

func main() {

	// CLI ARGS
	host := flag.String("host", "", "HOST (required) - Which host do you wish to proxy\nUsage Example:\n	$ go run proxy.go http://google.com/")
	port := flag.String("p", "1111", "PORT")
	id := flag.String("id", "", "ID (optional) - will set the default host id instead of the full domain name"+
		"\nUsage Example:\n	$ go run proxy.go randomnumberapi.com -id random -cache")
	cache := flag.Bool("cache", false, "CACHE (optional) - This will make proxy return from cache if possible "+
		"(from default [host].json unless -alt was set)\nUsage Example:\n	$ go run proxy.go http://google.com/ -cache")
	alt := flag.String("alt", "", "ALT (optional) [JSONFILE] - This will make proxy return from alternative cache file if possible"+
		"\nUsage Example:\n	$ go run proxy.go http://google.com/ -cache -alt ./mockMaps/google_alt.json		# respond from google_alt.json")
	strict := flag.Bool("strict", false, "STRICT (optional) - This will make proxy return ONLY from cache, no external calls")
	help := flag.Bool("h", false, "Shows this help message")
	noSave := flag.Bool("no-save", false, "NO-SAVE (optional) will not store any data from proxy")

	flag.Parse()
	if *help || (*host == "" && flag.NArg() == 0) {
		fmt.Println("\ngo run proxy.go [host] -p [port] OPTIONAL -cache -alt [JSONFILE] -strict\n")
		fmt.Println("	Usage Example:")
		fmt.Println("	$ go run proxy.go -host google.com/ -p 1111 -cache \n")
		flag.Usage()
	} else if *host == "" {
		if len(os.Args) > 0 {
			if os.Args[1] != "-host" {
				*host = os.Args[1]
				flag.CommandLine.Parse(append([]string{"-host"}, os.Args[1:]...))
			} else {
				*host = os.Args[1]
			}
		}
	}
	println()

	domain := getDomain(*host)
	if *id != "" {
		domain = *id
	} else {
		*id = domain
	}

	mockfile := mockFolder + domain + ".json"
	if *alt != "" {
		mockfile = mockFolder + *alt
	}

	if *host == "" {
		println("\n [host] is required. Exiting")
		os.Exit(1)
	}
	malicious := false // default

	startEpochUpdate(noSave)

	process := proxyProcess{
		id:        domain,
		port:      *port,
		host:      *host,
		mockfile:  mockfile,
		mock:      &mockMap{requests: map[string]string{}},
		malicious: false,
		cache:     *cache,
		strict:    *strict,
		noSave:    *noSave,
	}
	proxies = append(proxies, process)

	if !malicious {
		process.handler = process.LavaTestProxy
	} else {
		println()
		println("MMMMMMMMMMMMMMM MALICIOUS MMMMMMMMMMMMMMM PORT", port)
		println()
		// TODO: Make malicious proxy
		process.handler = process.LavaTestProxy
	}
	startProxyProcess(process)
}

func startProxyProcess(process proxyProcess) {
	process.mock.requests = jsonFileToMap(process.mockfile)
	if process.mock.requests == nil {
		process.mock.requests = map[string]string{}
	}
	if process.malicious {
		fakeResponse = true
	}
	fmt.Println(":::::::::::::::::::::::::::::::::::::::::::::::") // HOST ", process.host)
	fmt.Println("::::::::::::: Mock Proxy Started ::::::::::::::") // CACHE", fmt.Sprintf("%d", len(process.mock.requests)))
	fmt.Println(":::::::::::::::::::::::::::::::::::::::::::::::") // PORT ", process.port)
	println()
	fmt.Print(fmt.Sprintf(" ::: Proxy ID 		::: %s", process.id) + "\n")
	fmt.Print(fmt.Sprintf(" ::: Proxy Host 	::: %s", process.host) + "\n")
	fmt.Print(fmt.Sprintf(" ::: Return Cache 	::: %t", process.cache) + "\n")
	fmt.Print(fmt.Sprintf(" ::: Strict Mode 	::: %t", process.strict) + "\n")
	fmt.Print(fmt.Sprintf(" ::: Saving	 	::: %t", !process.noSave) + "\n")
	if !process.noSave || process.cache {
		fmt.Print(fmt.Sprintf(" ::: Cache File 	::: %s", process.mockfile) + "\n")
		fmt.Print(fmt.Sprintf(" ::: Loaded Responses 	::: %d", len(process.mock.requests)) + "\n")
	}
	println()
	fmt.Print(fmt.Sprintf(" ::: Proxy Started! 	::: ID: %s", process.id) + "\n")
	fmt.Print(fmt.Sprintf(" ::: Listening On 	::: %s", "http://0.0.0.0:"+process.port+"/") + "\n")

	http.HandleFunc("/", process.handler)
	err := http.ListenAndServe(":"+process.port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getMockBlockNumber() (block string) {
	return "0xe" + fmt.Sprintf("%d", epochCount) // return "0xe39ab8"
}

func startEpochUpdate(noSave *bool) {
	go func() {
		count := 0
		for {
			wait := 1 * time.Second
			time.Sleep(wait)
			if count%epochTime == 0 {
				epochCount += 1
			}
			if !*noSave && responsesChanged && count%saveJsonEvery == 0 {
				for _, process := range proxies {
					mapToJsonFile(*process.mock, process.mockfile)
				}
				responsesChanged = false
			}
			count += 1
		}
	}()
}

func fakeResult(val string, fake string) string {
	parts := strings.Split(val, ",")
	found := -1
	for i, part := range parts {
		if strings.Contains(part, "result") {
			found = i
		}
	}
	if found != -1 {
		parts[found] = fmt.Sprintf("\"result\":\"%s\"}", fake)
	}
	return strings.Join(parts, ",")
}
func (p proxyProcess) LavaTestProxy(rw http.ResponseWriter, req *http.Request) {

	host := p.host
	mock := p.mock

	// Get request body
	rawBody := getDataFromIORead(&req.Body, true)
	println()
	println(" ::: "+p.port+" ::: "+p.id+" ::: INCOMING PROXY MSG :::", string(rawBody))

	// TODO: make generic
	// Check if asking for blockNumber
	if fakeResponse && strings.Contains(string(rawBody), "blockNumber") {
		println("!!!!!!!!!!!!!! block number")
		rw.WriteHeader(200)
		rw.Write([]byte(fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"%s\"}", getMockBlockNumber())))

	} else {
		// Return Cached data if found in history and fromCache is set on
		if val, ok := mock.requests[string(rawBody)]; ok && p.cache {
			println(" ::: "+p.port+" ::: "+p.id+" ::: Cached Response ::: ", string(val))
			cacheCount += 1

			// Change Response
			if fakeResponse {
				val = fakeResult(val, "0xe000000000000000000")
				// val = "{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0xe000000000000000000\"}"
				println(p.port+" ::: Fake Response ::: ", val)
				fakeCount += 1
			}
			rw.WriteHeader(200)
			rw.Write([]byte(val))

		} else {
			// Recreating Request
			proxyRequest, err := createProxyRequest(req, host, string(rawBody))
			if err != nil {
				println(err.Error())
			} else {

				// Send Request to Host & Get Response
				proxyRes, err := sendRequest(proxyRequest)
				// respBody := []byte("error")
				var respBody []byte
				respBodyStr := "xxxxxx"
				status := 400
				if err != nil {
					println(err.Error())
					respBody = []byte(err.Error())
				} else {
					status = proxyRes.StatusCode
					respBody = getDataFromIORead(&proxyRes.Body, true)
					respBodyStr = string(respBody)
					mock.requests[string(rawBody)] = respBodyStr
					realCount += 1
					println(" ::: "+p.port+" ::: "+p.id+" ::: Real Response ::: ", respBodyStr)
				}

				// Check if response is not good, if not - try again
				if false && (strings.Contains(string(respBody), "error") || strings.Contains(string(respBody), "Error")) {
					println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX Got error in response - retrying request")

					// Recreating Request
					proxyRequest, err = createProxyRequest(req, host, string(rawBody))
					if err != nil {
						println(err.Error())
						respBody = []byte(err.Error())
					} else {

						// Send Request to Host & Get Response
						proxyRes, err = sendRequest(proxyRequest)
						if err != nil {
							println(err.Error())
							respBody = []byte("error: " + err.Error())
						} else {
							respBody = getDataFromIORead(&proxyRes.Body, true)
							mock.requests[string(rawBody)] = string(respBody)
							status = proxyRes.StatusCode
						}
						realCount += 1
						println(" ::: "+p.port+" ::: "+p.id+" ::: Real Response ::: ", string(respBody))

						// TODO: Check if response is good, if not - try again
						if strings.Contains(string(respBody), "error") || strings.Contains(string(respBody), "Error") {
							println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX Got another error in response ")
							println()
						} else {
							println("YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY SUCCESS - no error in response ")
							println()
						}
					}

				}

				// Change Response
				if fakeResponse {
					// respBody = []byte("{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0xe000000000000000000\"}")
					respBody = []byte(fakeResult(respBodyStr, "0xe000000000000000000"))
					println(" ::: "+p.port+" ::: "+p.id+" ::: Fake Response ::: ", string(respBody))
					fakeCount += 1
				}
				responsesChanged = true

				//Return Response
				if respBody == nil {
					respBody = []byte("error")
				}
				returnResponse(rw, status, respBody)
			}
		}
	}
	println("_________________________________", realCount, "/", cacheCount, "\n")
}