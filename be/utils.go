package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// FileValidator handles file validation operations
type FileValidator struct{}

// NewFileValidator creates a new file validator instance
func NewFileValidator() *FileValidator {
	return &FileValidator{}
}

// ValidateFile performs comprehensive validation on a file
func (v *FileValidator) ValidateFile(filePath string) error {
	// Check if file path is empty
	if strings.TrimSpace(filePath) == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	// Check if file exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}
	if err != nil {
		return fmt.Errorf("failed to access file: %w", err)
	}

	// Check if it's a directory
	if fileInfo.IsDir() {
		return fmt.Errorf("path points to a directory, not a file: %s", filePath)
	}

	// Check file extension
	if !v.isPDFFile(filePath) {
		ext := filepath.Ext(filePath)
		if ext == "" {
			return fmt.Errorf("file has no extension, expected .pdf")
		}
		return fmt.Errorf("invalid file extension '%s', expected .pdf", ext)
	}

	// Check file size (limit to 20MB)
	const maxFileSize = 20 * 1024 * 1024 // 20MB
	if fileInfo.Size() > maxFileSize {
		return fmt.Errorf("file too large: %s (max allowed: 20MB)", formatFileSize(fileInfo.Size()))
	}

	// Check if file is empty
	if fileInfo.Size() == 0 {
		return fmt.Errorf("file is empty")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// Read a small header to verify PDF signature and MIME type
	header := make([]byte, 1024)
	n, err := f.Read(header)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %w", err)
	}
	if n < 5 || !bytes.HasPrefix(header[:n], []byte("%PDF-")) {
		return fmt.Errorf("file does not appear to be a valid PDF (missing signature)")
	}

	mime := http.DetectContentType(header[:n])
	if mime != "application/pdf" && mime != "application/octet-stream" {
		return fmt.Errorf("invalid content type '%s', expected application/pdf", mime)
	}

	// Scan the first chunk for active content markers
	const sniffSize = 1024 * 1024 // 1MB
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to rewind file for scanning: %w", err)
	}

	body := make([]byte, sniffSize)
	m, _ := f.Read(body)
	lower := bytes.ToLower(body[:m])
	suspiciousMarkers := [][]byte{
		[]byte("/js"),
		[]byte("/javascript"),
		[]byte("/launch"),
		[]byte("/embeddedfile"),
		[]byte("/embeddedfiles"),
		[]byte("/richmedia"),
		[]byte("/openaction"),
		[]byte("/uri"),
	}
	for _, marker := range suspiciousMarkers {
		if bytes.Contains(lower, marker) {
			return fmt.Errorf("file contains disallowed PDF active content marker: %s", marker)
		}
	}

	return nil
}

// isPDFFile checks if the file has a PDF extension
func (v *FileValidator) isPDFFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".pdf"
}

// formatFileSize returns a human-readable file size
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// ErrorHandler provides centralized error handling
type ErrorHandler struct{}

// NewErrorHandler creates a new error handler instance
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// HandleError handles different types of errors with appropriate messages
func (h *ErrorHandler) HandleError(err error, context string) {
	if err == nil {
		return
	}

	fmt.Printf("Error in %s: %v\n", context, err)

	// Provide suggestions based on error type
	errorMsg := err.Error()
	switch {
	case strings.Contains(errorMsg, "does not exist"):
		fmt.Println("Suggestion: Check if the file path is correct and the file exists.")
	case strings.Contains(errorMsg, "permission denied"):
		fmt.Println("Suggestion: Check if you have permission to read the file.")
	case strings.Contains(errorMsg, "invalid file extension"):
		fmt.Println("Suggestion: Ensure the file has a .pdf extension.")
	case strings.Contains(errorMsg, "file too large"):
		fmt.Println("Suggestion: Try with a smaller PDF file (max 100MB).")
	case strings.Contains(errorMsg, "directory"):
		fmt.Println("Suggestion: Provide a path to a file, not a directory.")
	case strings.Contains(errorMsg, "no text could be extracted"):
		fmt.Println("Suggestion: The PDF might be image-based or encrypted. Try with a text-based PDF.")
	default:
		fmt.Println("Suggestion: Check the file and try again.")
	}
}
