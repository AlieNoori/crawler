package extractor

import "net/url"

type PageData struct {
	URL            string
	H1             string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func ExtractPageData(html, pageURL string) (*PageData, error) {
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}

	links, err := GetURLsFromHTML(html, baseURL)
	if err != nil {
		return nil, err
	}

	imgs, err := GetImagesFromHTML(html, baseURL)
	if err != nil {
		return nil, err
	}

	return &PageData{
		URL:            pageURL,
		H1:             GetH1FromHTML(html),
		FirstParagraph: GetFirstParagraphFromHTML(html),
		OutgoingLinks:  links,
		ImageURLs:      imgs,
	}, nil
}
