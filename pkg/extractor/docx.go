package extractor

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// DOCXExtractor extracts text from DOCX files by reading the embedded XML.
// DOCX files are ZIP archives containing word/document.xml with the text content.
type DOCXExtractor struct{}

func NewDOCXExtractor() *DOCXExtractor {
	return &DOCXExtractor{}
}

func (e *DOCXExtractor) SupportedFormats() []string {
	return []string{".docx"}
}

func (e *DOCXExtractor) Extract(filePath string) (string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return "", fmt.Errorf("docx open: %w", err)
	}
	defer r.Close()

	var documentXML *zip.File
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			documentXML = f
			break
		}
	}

	if documentXML == nil {
		return "", fmt.Errorf("word/document.xml not found in docx archive")
	}

	rc, err := documentXML.Open()
	if err != nil {
		return "", fmt.Errorf("docx read: %w", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return "", fmt.Errorf("docx readall: %w", err)
	}

	var doc document
	if err := xml.Unmarshal(data, &doc); err != nil {
		return "", fmt.Errorf("docx xml parse: %w", err)
	}

	var buf strings.Builder
	for _, para := range doc.Body.Paragraphs {
		for _, run := range para.Runs {
			buf.WriteString(run.Text)
		}
		buf.WriteRune('\n')
	}

	result := strings.TrimSpace(buf.String())
	if result == "" {
		return "", fmt.Errorf("no text content extracted from DOCX")
	}

	return result, nil
}

// XML structures for word/document.xml parsing
type document struct {
	XMLName xml.Name `xml:"document"`
	Body    body     `xml:"body"`
}

type body struct {
	Paragraphs []paragraph `xml:"p"`
}

type paragraph struct {
	Runs []run `xml:"r"`
}

type run struct {
	Text string `xml:"t"`
}
