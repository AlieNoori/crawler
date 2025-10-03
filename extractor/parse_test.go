package extractor

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestGetH1FromHTMLBasic(t *testing.T) {
	testCases := []struct {
		inputBody string
		expected  string
	}{
		{
			inputBody: "<html><body><h1>Test Title</h1></body></html>",
			expected:  "Test Title",
		},
		{
			inputBody: "<html><body><h1>   Test Title   </h1></body></html>",
			expected:  "Test Title",
		},
		{
			inputBody: "<html><body><h1>I like this</h1><h1>We dont like this one</h1></body></html>",
			expected:  "I like this",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test Case #%d", i+1), func(t *testing.T) {
			actual := GetH1FromHTML(tc.inputBody)

			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	testCases := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name: "first p tag of main tag",
			inputBody: `
<html>
	<body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body>
</html>`,
			expected: "Main paragraph.",
		},
		{
			name: "first p tag of document",
			inputBody: `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<h2>secondary header.</h2>
		</main>
	</body></html>`,
			expected: "Outside paragraph.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetFirstParagraphFromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	testCases := []struct {
		baseURL   string
		inputBody string
		expected  []string
	}{
		{
			baseURL:   "https://blog.boot.dev",
			inputBody: `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://blog.boot.dev"},
		},
		{
			baseURL:   "https://blog.boot.dev",
			inputBody: `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a><a href="/shop.html">Shop</a></body></html>`,
			expected:  []string{"https://blog.boot.dev", "https://blog.boot.dev/shop.html"},
		},
		{
			baseURL: "https://blog.boot.dev",
			inputBody: `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a><a href="/shop.html">Shop</a>
			<a href="https://peps.python.org/pep-0020/"><span>The Zen of Python</span></a></body></html>`,
			expected: []string{"https://blog.boot.dev", "https://blog.boot.dev/shop.html", "https://peps.python.org/pep-0020/"},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test Case #%d", i+1), func(t *testing.T) {
			baseURL, err := url.Parse(tc.baseURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := GetURLsFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestSetBaseURL(t *testing.T) {
	testCases := []struct {
		baseURL  string
		path     string
		expected string
	}{
		{
			baseURL:  "https://blog.boot.dev",
			path:     "/logo.png",
			expected: "https://blog.boot.dev/logo.png",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test #%d", i+1), func(t *testing.T) {
			baseURL, err := url.Parse(tc.baseURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual := SetBaseURL(tc.path, baseURL)

			if actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestHasBaseURL(t *testing.T) {
	testCases := []struct {
		name     string
		baseURL  string
		inputURL string
		expected bool
	}{
		{
			name:     "have baseURL",
			baseURL:  "https://blog.boot.dev",
			inputURL: "https://blog.boot.dev/path",
			expected: true,
		},
		{
			name:     "doesn't have baseURL",
			baseURL:  "https://blog.boot.dev",
			inputURL: "/logo.png",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.baseURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual := HasBaseURL(tc.inputURL, baseURL)
			if actual != tc.expected {
				t.Error("expected to be true")
			}
		})
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	testCases := []struct {
		baseURL   string
		inputBody string
		expected  []string
	}{
		{
			baseURL:   "https://blog.boot.dev",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://blog.boot.dev/logo.png"},
		},
		{
			baseURL:   "https://blog.boot.dev",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"/><img src="/bg.webp" alt="background" /></body></html>`,
			expected:  []string{"https://blog.boot.dev/logo.png", "https://blog.boot.dev/bg.webp"},
		},
		{
			baseURL:   "https://blog.boot.dev",
			inputBody: `<html><body><img src="https://blog.boot.dev/logo.png" alt="Logo"/><img src="/bg.webp" alt="background" /></body></html>`,
			expected:  []string{"https://blog.boot.dev/logo.png", "https://blog.boot.dev/bg.webp"},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test #%d", i+1), func(t *testing.T) {
			baseURL, err := url.Parse(tc.baseURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := GetImagesFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
