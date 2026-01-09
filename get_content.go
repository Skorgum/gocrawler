package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getH1FromHTML(html string) string {
	//extract H1 from HTML
	title := strings.NewReader(html)
	titlestr, _ := goquery.NewDocumentFromReader(title)
	h1 := titlestr.Find("h1").First().Text()
	if h1 != "" {
		return h1
	}
	return ""
}

func getFirstParagraphFromHTML(html string) string {
	//extract first paragraph from HTML, prioritizing <main> content
	paragraph := strings.NewReader(html)
	paragraphstr, _ := goquery.NewDocumentFromReader(paragraph)

	mainContent := paragraphstr.Find("main")
	if mainContent.Length() > 0 {
		p := mainContent.Find("p").First().Text()
		if p != "" {
			return p
		}
	}

	p := paragraphstr.Find("p").First().Text()
	if p != "" {
		return p
	}
	return ""
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}
		href = strings.TrimSpace(href)
		if href == "" {
			return
		}

		u, err := url.Parse(href)
		if err != nil {
			return
		}

		resolvedURL := baseURL.ResolveReference(u)
		urls = append(urls, resolvedURL.String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var images []string
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok {
			return
		}
		src = strings.TrimSpace(src)
		if src == "" {
			return
		}

		u, err := url.Parse(src)
		if err != nil {
			return
		}

		resolvedURL := baseURL.ResolveReference(u)
		images = append(images, resolvedURL.String())
	})

	return images, nil
}

func extractPageData(html, pageURL string) PageData {
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		log.Printf("could not parse pageURL %q: %v", pageURL, err)
		// you can still proceed with nil baseURL if needed
		parsedURL = nil
	}

	links, err := getURLsFromHTML(html, parsedURL)
	if err != nil {
		log.Printf("could not extract URLs from HTML: %v", err)
		links = []string{}
	}

	images, err := getImagesFromHTML(html, parsedURL)
	if err != nil {
		log.Printf("could not extract image URLs from HTML: %v", err)
		images = []string{}
	}

	return PageData{
		URL:            pageURL,
		H1:             getH1FromHTML(html),
		FirstParagraph: getFirstParagraphFromHTML(html),
		OutgoingLinks:  links,
		ImageURLs:      images,
	}
}
