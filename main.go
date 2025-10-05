package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/alieNoori/crawler/extractor"
)

type config struct {
	pages              map[string]*extractor.PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Print("no website provided\n")
		os.Exit(1)
	} else if len(args) > 3 {
		fmt.Print("too many arguments provided\n")
		os.Exit(1)
	}

	fmt.Printf("starting crawl\n-%s\n", args[0])

	rawURL := args[0]
	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalln(err)
	}
	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalln(err)
	}

	baseURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		pages:              map[string]*extractor.PageData{},
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	cfg.wg.Add(1)
	cfg.crawlPage(rawURL)
	cfg.wg.Wait()
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer cfg.wg.Done()

	defer func() {
		<-cfg.concurrencyControl
	}()

	if len(cfg.pages) >= cfg.maxPages {
		return
	}

	if !extractor.HasBaseURL(rawCurrentURL, cfg.baseURL) {
		return
	}

	normalizedURL, err := extractor.NormalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	cfg.mu.Lock()
	if _, ok := cfg.pages[normalizedURL]; ok {
		cfg.mu.Unlock()
		return
	}
	cfg.pages[normalizedURL] = nil
	cfg.mu.Unlock()

	htmlBody, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	pageData, err := extractor.ExtractPageData(htmlBody, cfg.baseURL.String())
	if err != nil {
		return
	}

	cfg.mu.Lock()
	cfg.pages[normalizedURL] = pageData
	cfg.mu.Unlock()

	if len(pageData.OutgoingLinks) <= 0 {
		return
	}

	for _, url := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
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
