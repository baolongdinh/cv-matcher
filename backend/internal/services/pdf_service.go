package services

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"encoding/hex"

	"github.com/ledongthuc/pdf"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) GenerateFileHash(pdfData []byte) string {
	hash := sha256.Sum256(pdfData)
	return hex.EncodeToString(hash[:])
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

	if extractedText.Len() == 0 {
		return "", numPages, fmt.Errorf("no text could be extracted from the PDF (might be scanned/image-based)")
	}

	return extractedText.String(), numPages, nil
}
