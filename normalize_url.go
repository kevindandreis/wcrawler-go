package main

import (
	"net/url"
	"strings"
)

func normalizeURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	parsedURL.Scheme = ""
	parsedURL.Host = strings.ToLower(parsedURL.Host)
	normalized := parsedURL.String()
	normalized = strings.TrimPrefix(normalized, "//")
	normalized = strings.TrimSuffix(normalized, "/")

	return normalized, nil
}