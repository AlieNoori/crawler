package main

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/alieNoori/crawler/extractor"
)

func writeCSVReport(pages map[string]*extractor.PageData, filename string) error {
	csvFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(csvFile)

	err = writer.Write([]string{"page_url", "h1", "first_paragraph", "outgoing_link_urls", "image_urls"})
	if err != nil {
		return err
	}

	for url, page := range pages {
		writer.Write(
			[]string{
				url,
				page.H1,
				page.FirstParagraph,
				strings.Join(page.OutgoingLinks, ";"),
				strings.Join(page.ImageURLs, ";"),
			})
	}

	return nil
}
