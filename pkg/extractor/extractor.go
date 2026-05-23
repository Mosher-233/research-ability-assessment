package extractor

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ContentExtractor defines the interface for extracting text from files.
type ContentExtractor interface {
	Extract(filePath string) (string, error)
	SupportedFormats() []string
}

// ExtractionResult holds the result of text extraction.
type ExtractionResult struct {
	Content    string  `json:"content"`
	Confidence float64 `json:"confidence"` // 0.0-1.0
	PageCount  int     `json:"page_count"`
	Metadata   string  `json:"metadata"` // human-readable extraction info
}

// ExtractorChain tries a sequence of extractors, returning the first successful result.
type ExtractorChain struct {
	extractors []ContentExtractor
}

// NewExtractorChain creates a chain with all available extractors.
func NewExtractorChain() *ExtractorChain {
	return &ExtractorChain{
		extractors: []ContentExtractor{
			NewPDFExtractor(),
			NewDOCXExtractor(),
		},
	}
}

// Extract tries each extractor based on file extension, then falls back to
// trying all extractors. Returns empty string if no extractor succeeds.
func (c *ExtractorChain) Extract(filePath string) *ExtractionResult {
	ext := strings.ToLower(filepath.Ext(filePath))

	// First, try extractors that claim to support this extension
	for _, e := range c.extractors {
		for _, f := range e.SupportedFormats() {
			if strings.EqualFold(ext, f) {
				content, err := e.Extract(filePath)
				if err == nil && strings.TrimSpace(content) != "" {
					return &ExtractionResult{
						Content:    content,
						Confidence: 0.95,
						Metadata:   fmt.Sprintf("extracted via %T", e),
					}
				}
			}
		}
	}

	// Fallback: try all extractors regardless of extension
	for _, e := range c.extractors {
		content, err := e.Extract(filePath)
		if err == nil && strings.TrimSpace(content) != "" {
			return &ExtractionResult{
				Content:    content,
				Confidence: 0.7,
				Metadata:   fmt.Sprintf("fallback extraction via %T", e),
			}
		}
	}

	return nil
}

// GetSupportedFormats returns all supported extensions across all extractors.
func (c *ExtractorChain) GetSupportedFormats() []string {
	seen := make(map[string]bool)
	var formats []string
	for _, e := range c.extractors {
		for _, f := range e.SupportedFormats() {
			if !seen[f] {
				seen[f] = true
				formats = append(formats, f)
			}
		}
	}
	return formats
}
