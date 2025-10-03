package extractor

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetH1FromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	h1 := strings.TrimSpace(doc.Find("h1").First().Text())

	return h1
}

func GetFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	var firstParagraph string

	firstParagraph = doc.Find("p").First().Text()

	main := doc.Find("main")
	if main.Length() > 0 {
		if main.Find("p").Length() > 0 {
			firstParagraph = main.Find("p").First().Text()
		}
	}

	return strings.TrimSpace(firstParagraph)
}

func GetURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	var urls []string

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if exists {
			if HasBaseURL(url, baseURL) {
				urls = append(urls, url)
			} else {
				urls = append(urls, SetBaseURL(url, baseURL))
			}
		}
	})

	return urls, nil
}

func HasBaseURL(url string, baseURL *url.URL) bool {
	return strings.HasPrefix(url, fmt.Sprint(baseURL))
}

func SetBaseURL(url string, baseURL *url.URL) string {
	if strings.HasPrefix(url, "http") {
		return url
	} else {
		rawBaseURL, _ := strings.CutSuffix(baseURL.String(), "/")
		return rawBaseURL + url
	}
}

func GetImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	var urls []string

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		url, exists := s.Attr("src")
		if exists {
			if HasBaseURL(url, baseURL) {
				urls = append(urls, url)
			} else {
				urls = append(urls, SetBaseURL(url, baseURL))
			}
		}
	})

	return urls, nil
}
