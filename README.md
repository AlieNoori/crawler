# Web Crawler

A concurrent web crawler built in Go that extracts page data and generates CSV reports. This crawler traverses websites, extracts key information from HTML pages, and compiles the results into a structured CSV format.

## Features

- **Concurrent Crawling**: Configurable concurrency control for efficient parallel page processing
- **Comprehensive Data Extraction**: Extracts H1 headers, first paragraphs, outgoing links, and image URLs
- **URL Normalization**: Handles relative and absolute URLs with proper normalization
- **CSV Report Generation**: Outputs structured data in CSV format for easy analysis
- **Configurable Limits**: Set maximum pages to crawl and concurrency levels
- **Same-Domain Crawling**: Automatically stays within the base domain

## Installation

```bash
git clone https://github.com/alieNoori/crawler.git
cd crawler
go mod download
```

## Usage

```bash
go run . <base_url> <max_concurrency> <max_pages>
```

### Parameters

- `base_url`: The starting URL for the crawl (e.g., `https://example.com`)
- `max_concurrency`: Maximum number of concurrent requests (e.g., `5`)
- `max_pages`: Maximum number of pages to crawl (e.g., `100`)

### Example

```bash
go run . https://blog.boot.dev 10 50
```

This command will:
- Start crawling from `https://blog.boot.dev`
- Use up to 10 concurrent goroutines
- Stop after crawling 50 pages

## Output

The crawler generates a `report.csv` file with the following columns:

| Column | Description |
|--------|-------------|
| `page_url` | The URL of the crawled page |
| `h1` | The first H1 header found on the page |
| `first_paragraph` | The first paragraph text (prioritizes `<main>` content) |
| `outgoing_link_urls` | Semicolon-separated list of outgoing links |
| `image_urls` | Semicolon-separated list of image URLs |

## Project Structure

```
.
├── main.go                    # Main crawler logic and configuration
├── csv_report.go              # CSV report generation
├── extractor/
│   ├── extract.go            # Main extraction orchestrator
│   ├── extract_test.go       # Integration tests
│   ├── parse.go              # HTML parsing functions
│   ├── parse_test.go         # Parser unit tests
│   ├── normalize_url.go      # URL normalization
│   └── normalize_url_test.go # Normalization tests
├── go.mod                     # Go module definition
└── .gitignore                # Git ignore rules
```

## How It Works

1. **Initialization**: The crawler starts with a base URL and creates a configuration with concurrency controls
2. **Page Crawling**: For each page:
   - Fetches the HTML content
   - Extracts H1, first paragraph, links, and images
   - Normalizes all URLs
   - Stores page data in memory
3. **Recursive Crawling**: Follows outgoing links that belong to the same domain
4. **Concurrency Control**: Uses channels and wait groups to manage concurrent requests
5. **Report Generation**: Writes all collected data to a CSV file

## Key Components

### Extractor Package

- **`ExtractPageData`**: Main function that orchestrates all extraction operations
- **`GetH1FromHTML`**: Extracts the first H1 header from HTML
- **`GetFirstParagraphFromHTML`**: Extracts the first paragraph (prioritizes `<main>` tag content)
- **`GetURLsFromHTML`**: Extracts and normalizes all anchor links
- **`GetImagesFromHTML`**: Extracts and normalizes all image URLs
- **`NormalizeURL`**: Strips protocol and trailing slashes for URL deduplication

### Concurrency Design

The crawler uses several Go concurrency patterns:
- **Mutex**: Protects shared page data map
- **Buffered Channel**: Controls maximum concurrent requests
- **WaitGroup**: Ensures all goroutines complete before generating report

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run specific package tests:

```bash
go test ./extractor
```

## Dependencies

- [goquery](https://github.com/PuerkitoBio/goquery) - HTML parsing and DOM manipulation

## Limitations

- Only crawls pages within the same domain as the base URL
- Does not execute JavaScript (static HTML only)
- Respects no robots.txt or meta robots tags
- No rate limiting (be mindful of server load)

## Future Enhancements

- [ ] Add robots.txt support
- [ ] Implement politeness delays
- [ ] Support for JavaScript-rendered content
- [ ] Add progress indicators
- [ ] Export to additional formats (JSON, SQLite)
- [ ] Configurable extraction rules
- [ ] Resume interrupted crawls

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
