package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: parsing current URL: %v\n", err)
		return
	}

	if cfg.baseURL.Host != currentURL.Host {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: normalizing URL: %v\n", err)
		return
	}

	cfg.mu.Lock()
	if _, ok := cfg.pages[normalizedURL]; ok {
		cfg.mu.Unlock()
		return
	}

	// For now, satisfy the map structure while we perform the work.
	// We'll update the value once we've extracted the data.
	cfg.pages[normalizedURL] = PageData{}
	cfg.mu.Unlock()

	fmt.Printf("crawling %s\n", rawCurrentURL)
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: fetching HTML: %v\n", err)
		return
	}

	data := extractPageData(html, rawCurrentURL)

	cfg.mu.Lock()
	cfg.pages[normalizedURL] = data
	cfg.mu.Unlock()

	for _, nextURL := range data.OutgoingLinks {
		cfg.crawlPage(nextURL)
	}
}
