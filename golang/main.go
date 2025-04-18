// Package implements collection http endpoints availability based on provided config file.
// The app expects a file with configuration, structured in yaml format.
//
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Endpoint is the representation of a configuration for each section of config.
type Endpoint struct {
	Name    string            `yaml:"name"`
	URL     string            `yaml:"url"`
	Method  string            `yaml:"method"`
	Headers map[string]string `yaml:"headers"`
	Body    string            `yaml:"body"`
}

//DomainStats represents counters of attempts to reach remote URL.
type DomainStats struct {
	Success int
	Total   int
}

// various variables initialization of the app
var (
	stats    = make(map[string]*DomainStats)
	infoLog  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// checkHealth function takes one endpoint struct and makes a http call to the configured url and increments stats' counters.
// it receives the instance of endpoint  to process
// it returns nothing
// executed as a gorouting
func checkHealth(endpoint Endpoint) {
	debugLogFunc("goroutine checkHealth for: " + endpoint.Name)
	domain := extractDomain(endpoint.URL)
	debugLogFunc("   stat success/total=" + fmt.Sprint(stats[domain].Success) + "/" + fmt.Sprint(stats[domain].Total))

	var client = &http.Client{}

	bodyBytes, err := json.Marshal(endpoint.Body)
	if err != nil {
		return
	}
	//commented some lines as too deep information for daily troubleshooting, was used for feature checking
	//debugLogFunc("    bodyBytes= " + fmt.Sprint(bodyBytes))
	reqBody := bytes.NewReader(bodyBytes)
	//debugLogFunc("    reqBody= " + fmt.Sprint(reqBody))

	req, err := http.NewRequest(endpoint.Method, endpoint.URL, reqBody)
	if err != nil {
		errorLog.Println("Error creating request:", err)
		return
	}
	//debugLogFunc("    req= " + fmt.Sprint(req))

	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
		//debugLogFunc("    headers: key" + fmt.Sprint(key) + " value=" + fmt.Sprint(value))
	}
	reqstart := time.Now()
	resp, err := client.Do(req)
	reqcompleted := time.Since(reqstart).Milliseconds()
	stats[domain].Total++
	if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 && reqcompleted < 500 {
		stats[domain].Success++
		debugLogFunc("   request \"" + endpoint.Name + "\" for " + domain + " is succedded and added to the success counter")
	} else {
		if err == nil {
			debugLogFunc("   request to " + domain + " failed due to condtions, code=" + fmt.Sprint(resp.StatusCode) + " in " + fmt.Sprint(reqcompleted) + "ms")
		} else {
			errorLog.Printf("    request failed with " + fmt.Sprint(err))
		}
	}
	debugLogFunc("   domain=" + domain)

}

// extractDomain function trim scheme and port for given url
// takes url string
// return domain string
// used as a helper for collecting stats
func extractDomain(url string) string {
	//cut scheme
	urlSplit := strings.Split(url, "//")
	domain := strings.Split(urlSplit[len(urlSplit)-1], "/")[0]
	//cut port
	urlNoPort := strings.Split(domain, ":")
	domain = urlNoPort[0]
	return domain
}

// monitorEndpoints is a function with main loop which is executing goroutines for each config's entry
// takes list of endpoints from config file
// return nothing
// used as a main cycle for the script
func monitorEndpoints(endpoints []Endpoint) {
	for _, endpoint := range endpoints {
		debugLogFunc("Initializing over=" + fmt.Sprint(endpoint.Name))
		debugLogFunc("  Endpoint=" + fmt.Sprint(endpoint.URL))
		domain := extractDomain(endpoint.URL)
		debugLogFunc("  Domain=" + fmt.Sprint(domain))
		if stats[domain] == nil {
			stats[domain] = &DomainStats{}
		}
	}

	debugLogFunc("Main checking loop")
	for {
		for _, endpoint := range endpoints {
			go checkHealth(endpoint)
		}
		time.Sleep(15 * time.Second)
	}
}

// logTimer function is goroutine which show stats, indepently from main loop in monitorEndpoints, with it's own timer
// takes nothing
// return nothing
// used as an independent logs renderer
func logTimer() {
	for {
		time.Sleep(15 * time.Second)
		logResults()
	}
}

// logResults function is preparing stats metrics to log formatted messages and calculate availability percentage based on stats
// takes nothing
// return nothing
func logResults() {
	debugLogFunc("gorutine for logResults started")
	for domain, stat := range stats {
		percentage := int(math.Round(100 * float64(stat.Success) / float64(stat.Total)))
		debugLogFunc(domain + " success/total=" + fmt.Sprint(stat.Success) + "/" + fmt.Sprint(stat.Total) + "=" + fmt.Sprint(percentage))
		infoLog.Printf("%s has %d%% availability\n", domain, percentage)
	}
}

// debugLogFunc function is sending message to logs if debug mode is enabled only
// takes message string and envvar LOG_LEVEL
// return nothing
func debugLogFunc(message string) {
	lvl, _ := os.LookupEnv("LOG_LEVEL")
	if lvl == "DEBUG" {
		debugLog.Print(message)
	}
}

// validateEndpoints function checks if config entry has all required fields to check availabiliry
// takes reference of endpoints from config file
// returns new slice with filtered etries only
func validateEndpoints(e *[]Endpoint) []Endpoint {
	var outEndoints []Endpoint
	for _, endpoint := range *e {
		urlExist, nameExist := "", ""
		urlExist, nameExist = endpoint.URL, endpoint.Name
		if urlExist == "" || nameExist == "" {
			errorLog.Printf("Validation failed: URL or Name not exist, skipping...")
		} else {
			debugLogFunc("Validation passed url=" + urlExist + " name=" + nameExist)
			outEndoints = append(outEndoints, endpoint)
		}
	}
	debugLogFunc("Num of entries in config file: " + fmt.Sprint(len(*e)))
	debugLogFunc("Num of entries after validation: " + fmt.Sprint(len(outEndoints)))
	return outEndoints
}

// main function is an app entrypoint
// it validates input file and unmarchalling to structure
// takes argument as a path to config file
// return nothing
func main() {

	if len(os.Args) < 2 {
		errorLog.Fatal("Usage: go run main.go <config_file>")
	}

	filePath := os.Args[1]
	debugLogFunc("Reading config file" + filePath)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		errorLog.Fatal("Error reading file:", err)
	}

	var endpoints []Endpoint
	if err := yaml.Unmarshal(data, &endpoints); err != nil {
		errorLog.Fatal("Error parsing YAML:", err)
	}
	filteredEndoints := validateEndpoints(&endpoints)
	debugLogFunc("Checking endpoints")
	go logTimer()
	monitorEndpoints(filteredEndoints)
}
