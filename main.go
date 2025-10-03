package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/alieNoori/crawler/extractor"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Print("no website provided\n")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Print("too many arguments provided\n")
		os.Exit(1)
	}

	fmt.Printf("starting crawl\n-%s\n", args[0])

	pages := map[string]int{}

	crawlPage(args[0], args[0], pages)
}

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "BootCrawler/1.0")

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return "", nil
	}
	defer response.Body.Close()

	if response.StatusCode > 400 {
		return "", errors.New("client level error")
	}

	// if contentType := response.Header.Get("Content-Type"); contentType != "" && contentType != "text/html" {
	// 	return "", errors.New("website doesn't return text/html type")
	// }

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return
	}

	if !extractor.HasBaseURL(rawCurrentURL, baseURL) {
		return
	}

	normalizedURL, err := extractor.NormalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if _, ok := pages[normalizedURL]; ok {
		pages[normalizedURL]++
		return
	} else {
		pages[normalizedURL] = 1
	}

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	pageData, err := extractor.ExtractPageData(htmlBody, rawBaseURL)
	if err != nil {
		return
	}

	if len(pageData.OutgoingLinks) <= 0 {
		return
	}

	fmt.Println("====================================")
	fmt.Println(rawCurrentURL)
	fmt.Println("------------------------------------")
	fmt.Printf("%+v\n", pageData)
	fmt.Println("====================================")

	for _, url := range pageData.OutgoingLinks {
		crawlPage(rawBaseURL, url, pages)
	}
}
