package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}
	h1 := doc.Find("h1, h2").First().Text()
	return strings.TrimSpace(h1)
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	main := doc.Find("main")
	var p string
	if main.Length() > 0 {
		p = main.Find("p").First().Text()
	} else {
		p = doc.Find("p").First().Text()
	}

	return strings.TrimSpace(p)
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	urls := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}

		u, err := url.Parse(href)
		if err != nil {
			return
		}

		resolved := baseURL.ResolveReference(u)
		urls = append(urls, resolved.String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	urls := []string{}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok {
			return
		}

		u, err := url.Parse(src)
		if err != nil {
			return
		}

		resolved := baseURL.ResolveReference(u)
		urls = append(urls, resolved.String())
	})

	return urls, nil
}

func extractPageData(html, pageURL string) PageData {
	baseURL, _ := url.Parse(pageURL)

	heading := getHeadingFromHTML(html)
	paragraph := getFirstParagraphFromHTML(html)
	urls, _ := getURLsFromHTML(html, baseURL)
	images, _ := getImagesFromHTML(html, baseURL)

	return PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: paragraph,
		OutgoingLinks:  urls,
		ImageURLs:      images,
	}
}
