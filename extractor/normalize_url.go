package extractor

import (
	"net/url"
	"strings"
)

func NormalizeURL(rawURL string) (string, error) {
	URL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	normalizedURL, _ := strings.CutSuffix(URL.Host+URL.Path, "/")

	return normalizedURL, nil
}
