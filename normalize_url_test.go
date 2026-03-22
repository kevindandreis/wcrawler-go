package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "remove scheme",
			input:    "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove trailing slash",
			input:    "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "lowercase host",
			input:    "https://BLOG.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove http",
			input:    "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGetHeadingFromHTML(t *testing.T) {
	cases := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name:      "basic h1",
			inputBody: "<html><body><h1>Test Title</h1></body></html>",
			expected:  "Test Title",
		},
		{
			name:      "h2 fallback",
			inputBody: "<html><body><h2>Fallback Title</h2></body></html>",
			expected:  "Fallback Title",
		},
		{
			name:      "no heading",
			inputBody: "<html><body><p>No heading</p></body></html>",
			expected:  "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getHeadingFromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	cases := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name:      "main priority",
			inputBody: `<html><body><p>Outside</p><main><p>Inside</p></main></body></html>`,
			expected:  "Inside",
		},
		{
			name:      "fallback to first p",
			inputBody: `<html><body><p>First</p><p>Second</p></body></html>`,
			expected:  "First",
		},
		{
			name:      "no paragraph",
			inputBody: `<html><body><div>No p</div></body></html>`,
			expected:  "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	cases := []struct {
		name         string
		inputBody    string
		inputBaseURL string
		expected     []string
	}{
		{
			name:         "absolute and relative URLs",
			inputBaseURL: "https://blog.boot.dev",
			inputBody:    `<html><body><a href="/path/one"><span>Boot.dev</span></a><a href="https://other.com/path/one"><span>Boot.dev</span></a></body></html>`,
			expected:     []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:         "no URLs",
			inputBaseURL: "https://blog.boot.dev",
			inputBody:    `<html><body><p>No links here</p></body></html>`,
			expected:     []string{},
		},
		{
			name:         "invalid URLs are skipped",
			inputBaseURL: "https://blog.boot.dev",
			inputBody:    `<html><body><a href=" :invalid">Link</a><a href="/valid">Link</a></body></html>`,
			expected:     []string{"https://blog.boot.dev/valid"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, _ := url.Parse(tc.inputBaseURL)
			actual, err := getURLsFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	cases := []struct {
		name         string
		inputBody    string
		inputBaseURL string
		expected     []string
	}{
		{
			name:         "absolute and relative image URLs",
			inputBaseURL: "https://blog.boot.dev",
			inputBody:    `<html><body><img src="/img1.png"><img src="https://other.com/img2.jpg"></body></html>`,
			expected:     []string{"https://blog.boot.dev/img1.png", "https://other.com/img2.jpg"},
		},
		{
			name:         "no images",
			inputBaseURL: "https://blog.boot.dev",
			inputBody:    `<html><body><p>No images here</p></body></html>`,
			expected:     []string{},
		},
		{
			name:         "skips img without src",
			inputBaseURL: "https://blog.boot.dev",
			inputBody:    `<html><body><img alt="no src"><img src="/valid.png"></body></html>`,
			expected:     []string{"https://blog.boot.dev/valid.png"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, _ := url.Parse(tc.inputBaseURL)
			actual, err := getImagesFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestExtractPageData(t *testing.T) {
	cases := []struct {
		name         string
		inputHTML    string
		inputPageURL string
		expected     PageData
	}{
		{
			name:         "full page data",
			inputPageURL: "https://blog.boot.dev/path/",
			inputHTML: `<html><body>
				<h1>The Heading</h1>
				<main><p>The paragraph.</p></main>
				<a href="/link1">Link</a>
				<img src="/img1.png">
			</body></html>`,
			expected: PageData{
				URL:            "https://blog.boot.dev/path/",
				Heading:        "The Heading",
				FirstParagraph: "The paragraph.",
				OutgoingLinks:  []string{"https://blog.boot.dev/link1"},
				ImageURLs:      []string{"https://blog.boot.dev/img1.png"},
			},
		},
		{
			name:         "no special tags, empty lists",
			inputPageURL: "https://blog.boot.dev",
			inputHTML:    `<html><body><div>Only a div</div></body></html>`,
			expected: PageData{
				URL:            "https://blog.boot.dev",
				Heading:        "",
				FirstParagraph: "",
				OutgoingLinks:  []string{},
				ImageURLs:      []string{},
			},
		},
		{
			name:         "h2 fallback and p outside main",
			inputPageURL: "https://blog.boot.dev",
			inputHTML: `<html><body>
				<h2>H2 Heading</h2>
				<p>Paragraph outside.</p>
			</body></html>`,
			expected: PageData{
				URL:            "https://blog.boot.dev",
				Heading:        "H2 Heading",
				FirstParagraph: "Paragraph outside.",
				OutgoingLinks:  []string{},
				ImageURLs:      []string{},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractPageData(tc.inputHTML, tc.inputPageURL)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}
