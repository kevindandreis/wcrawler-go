package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}

	rawBaseURL := os.Args[1]
	maxConcurrencyStr := os.Args[2]
	maxPagesStr := os.Args[3]

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Error - parsing base URL: %v\n", err)
		os.Exit(1)
	}

	maxConcurrency, err := strconv.Atoi(maxConcurrencyStr)
	if err != nil {
		fmt.Printf("Error - parsing maxConcurrency: %v\n", err)
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(maxPagesStr)
	if err != nil {
		fmt.Printf("Error - parsing maxPages: %v\n", err)
		os.Exit(1)
	}

	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	fmt.Printf("starting crawl of: %s\n", rawBaseURL)
	cfg.crawlPage(rawBaseURL)

	for url, data := range cfg.pages {
		fmt.Printf("%s: %d urls\n", url, len(data.OutgoingLinks))
	}

	err = writeJSONReport(cfg.pages, "report.json")
	if err != nil {
		fmt.Printf("Error - writeJSONReport: %v\n", err)
	}
}
