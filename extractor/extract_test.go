package extractor

import (
	"fmt"
	"reflect"
	"testing"
)

func TestExtractPageData(t *testing.T) {
	testCases := []struct {
		inputURL  string
		inputBody string
		expected  *PageData
	}{
		{
			inputURL: "https://blog.boot.dev",
			inputBody: `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`,
			expected: &PageData{
				URL:            "https://blog.boot.dev",
				H1:             "Test Title",
				FirstParagraph: "This is the first paragraph.",
				OutgoingLinks:  []string{"https://blog.boot.dev/link1"},
				ImageURLs:      []string{"https://blog.boot.dev/image1.jpg"},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test Case #%d", i+1), func(t *testing.T) {
			actual, err := ExtractPageData(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}
