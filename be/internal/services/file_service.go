package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

const (
	maxPDFSize   = int64(20 << 20)   // 20MB
	sniffPDFSize = int64(1 << 20)    // 1MB
	minPDFHeader = 5                 // "%PDF-" is 5 bytes
)

// ProcessUploadedFile extracts text content from uploaded files
func (fs *FileService) ProcessUploadedFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))

	if ext != ".pdf" {
		return "", fmt.Errorf("unsupported file type: %s (only PDF allowed)", ext)
	}

	content, err := readLimited(file, maxPDFSize)
	if err != nil {
		return "", err
	}

	if err := validatePDFContent(content, header.Filename); err != nil {
		return "", err
	}

	return fs.processPDFBuffer(content)
}

func (fs *FileService) processPDFBuffer(content []byte) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	var text strings.Builder
	numPages := reader.NumPage()

	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}

		// Try GetPlainText first
		pageText, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		// Clean and normalize the text
		cleanedText := fs.cleanPDFText(pageText)
		
		text.WriteString(cleanedText)
		text.WriteString("\n\n") // Add paragraph breaks between pages
	}

	if text.Len() == 0 {
		return "", fmt.Errorf("no text content found in PDF")
	}

	finalText := text.String()
	
	// Final cleaning pass
	finalText = fs.normalizePDFText(finalText)
	
	return finalText, nil
}

// cleanPDFText cleans up PDF text extraction artifacts
func (fs *FileService) cleanPDFText(text string) string {
	// First, remove common PDF artifacts
	text = strings.ReplaceAll(text, "□", " ")
	text = strings.ReplaceAll(text, "�", "")
	text = strings.ReplaceAll(text, "\u00a0", " ") // non-breaking space
	
	// Replace multiple spaces/newlines with single space
	text = strings.Join(strings.Fields(text), " ")
	
	// AGGRESSIVE word boundary detection
	var result strings.Builder
	runes := []rune(text)
	
	for i := 0; i < len(runes); i++ {
		current := runes[i]
		result.WriteRune(current)
		
		if i < len(runes)-1 {
			next := runes[i+1]
			
			// Skip if already has space
			if current == ' ' || next == ' ' {
				continue
			}
			
			// Add space between lowercase and uppercase (camelCase)
			// Example: "matriksselisih" → "matriks selisih"
			if (current >= 'a' && current <= 'z') && (next >= 'A' && next <= 'Z') {
				result.WriteRune(' ')
				continue
			}
			
			// Add space between uppercase and lowercase (if previous was lowercase)
			// Example: "SSdan" → "SS dan"
			if i > 0 && (current >= 'A' && current <= 'Z') && (next >= 'a' && next <= 'z') {
				prev := runes[i-1]
				if prev >= 'a' && prev <= 'z' {
					result.WriteRune(' ')
					continue
				}
			}
			
			// Add space between letter and number
			if ((current >= 'a' && current <= 'z') || (current >= 'A' && current <= 'Z')) && 
			   (next >= '0' && next <= '9') {
				result.WriteRune(' ')
				continue
			}
			
			// Add space between number and letter
			if (current >= '0' && current <= '9') && 
			   ((next >= 'a' && next <= 'z') || (next >= 'A' && next <= 'Z')) {
				result.WriteRune(' ')
				continue
			}
			
			// Add space after closing bracket if followed by letter
			if (current == ')' || current == ']' || current == '}') && 
			   ((next >= 'A' && next <= 'Z') || (next >= 'a' && next <= 'z')) {
				result.WriteRune(' ')
				continue
			}
			
			// Add space before opening bracket if preceded by letter
			if ((current >= 'a' && current <= 'z') || (current >= 'A' && current <= 'Z')) &&
			   (next == '(' || next == '[' || next == '{') {
				result.WriteRune(' ')
				continue
			}
			
			// Add space after period if followed by uppercase (sentence boundary)
			if current == '.' && next >= 'A' && next <= 'Z' {
				result.WriteRune(' ')
				continue
			}
			
			// Add space after comma if not already present
			if current == ',' && ((next >= 'A' && next <= 'Z') || (next >= 'a' && next <= 'z')) {
				result.WriteRune(' ')
				continue
			}
			
			// Add space after colon/semicolon if followed by letter
			if (current == ':' || current == ';') && 
			   ((next >= 'A' && next <= 'Z') || (next >= 'a' && next <= 'z')) {
				result.WriteRune(' ')
				continue
			}
		}
	}
	
	return result.String()
}

// normalizePDFText performs final normalization
func (fs *FileService) normalizePDFText(text string) string {
	// Remove excessive whitespace
	text = strings.Join(strings.Fields(text), " ")
	
	// Fix common ligatures and special characters
	replacements := map[string]string{
		"ﬁ": "fi",
		"ﬂ": "fl",
		"ﬀ": "ff",
		"ﬃ": "ffi",
		"ﬄ": "ffl",
		"□": " ", // Replace box character with space
		"�": "",  // Remove replacement character
	}
	
	for old, new := range replacements {
		text = strings.ReplaceAll(text, old, new)
	}
	
	// Ensure proper sentence spacing
	text = strings.ReplaceAll(text, ". ", ". ")
	text = strings.ReplaceAll(text, "? ", "? ")
	text = strings.ReplaceAll(text, "! ", "! ")
	
	// Remove multiple consecutive spaces again
	text = strings.Join(strings.Fields(text), " ")

	return strings.TrimSpace(text)
}

func readLimited(r io.Reader, limit int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("file too large (max allowed: %d bytes)", limit)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("file is empty")
	}
	return data, nil
}

func validatePDFContent(content []byte, filename string) error {
	if len(content) < minPDFHeader || !bytes.HasPrefix(content, []byte("%PDF-")) {
		return fmt.Errorf("file does not appear to be a valid PDF: %s", filename)
	}

	// Verify content type
	mime := http.DetectContentType(content[:min(len(content), 512)])
	if mime != "application/pdf" && mime != "application/octet-stream" {
		return fmt.Errorf("invalid content type '%s', expected application/pdf", mime)
	}

	// Scan for common active/suspicious markers
	end := int64(len(content))
	if end > sniffPDFSize {
		end = sniffPDFSize
	}
	lower := bytes.ToLower(content[:end])
	markers := [][]byte{
		[]byte("/js"),
		[]byte("/javascript"),
		[]byte("/launch"),
		[]byte("/embeddedfile"),
		[]byte("/embeddedfiles"),
		[]byte("/richmedia"),
		[]byte("/openaction"),
		[]byte("/uri"),
	}
	for _, marker := range markers {
		if bytes.Contains(lower, marker) {
			return fmt.Errorf("file contains disallowed PDF active content marker: %s", marker)
		}
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
