package services

import (
	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) ExtractText(pdfData []byte) (string, int, error) {
	reader := bytes.NewReader(pdfData)
	contentLength := int64(len(pdfData))
	
	r, err := pdf.NewReader(reader, contentLength)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create pdf reader: %v", err)
	}

	var extractedText bytes.Buffer
	numPages := r.NumPage()

	for i := 1; i <= numPages; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			// Skip problematic pages instead of failing entirely
			continue
		}

		extractedText.WriteString(text)
		extractedText.WriteString("\n")
	}

	return extractedText.String(), numPages, nil
}
