package extractor

import (
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

// PDFExtractor extracts text from PDF files.
type PDFExtractor struct{}

func NewPDFExtractor() *PDFExtractor {
	return &PDFExtractor{}
}

func (e *PDFExtractor) SupportedFormats() []string {
	return []string{".pdf"}
}

func (e *PDFExtractor) Extract(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("pdf open: %w", err)
	}
	defer f.Close()

	var buf strings.Builder
	totalPage := r.NumPage()

	for pageNum := 1; pageNum <= totalPage; pageNum++ {
		page := r.Page(pageNum)
		if page.V.IsNull() {
			continue
		}

		texts := page.Content().Text
		var lastY float64
		var lastChar rune

		for _, text := range texts {
			// Detect line breaks based on Y position changes
			if lastY != 0 && absDiff(text.Y, lastY) > 3 {
				if lastChar != '\n' && lastChar != 0 {
					buf.WriteRune('\n')
					lastChar = '\n'
				}
			}

			s := text.S
			for _, r := range s {
				buf.WriteRune(r)
				lastChar = r
			}
			lastY = text.Y
		}

		// Add page separator unless it's the last page
		if pageNum < totalPage {
			buf.WriteRune('\n')
		}
	}

	result := cleanPDFText(buf.String())
	if strings.TrimSpace(result) == "" {
		return "", fmt.Errorf("no text content extracted from PDF")
	}

	return result, nil
}

func absDiff(a, b float64) float64 {
	if a > b {
		return a - b
	}
	return b - a
}

// cleanPDFText normalizes whitespace in extracted PDF text.
func cleanPDFText(s string) string {
	var b strings.Builder
	prevSpace := false

	for _, r := range s {
		switch {
		case r == '\r':
			continue
		case r == '\n':
			b.WriteRune('\n')
			prevSpace = false
		case r == ' ' || r == '\t':
			if !prevSpace {
				b.WriteRune(' ')
				prevSpace = true
			}
		default:
			b.WriteRune(r)
			prevSpace = false
		}
	}

	// Collapse multiple consecutive newlines
	result := b.String()
	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}

	return strings.TrimSpace(result)
}
