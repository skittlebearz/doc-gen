package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type PDFRequest struct {
	HTML  string `json:"html"`
	Title string `json:"title,omitempty"`
}

type PDFResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	PDF     []byte `json:"pdf,omitempty"`
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/convert", convertHandler)

	log.Println("PDF Service starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PDFRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.HTML == "" {
		http.Error(w, "HTML content is required", http.StatusBadRequest)
		return
	}

	// Generate PDF using headless browser
	pdfData, err := generatePDF(req.HTML, req.Title)
	if err != nil {
		log.Printf("PDF generation failed: %v", err)
		http.Error(w, fmt.Sprintf("PDF generation failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Return PDF as response
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=document.pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfData)
}

func generatePDF(htmlContent, title string) ([]byte, error) {
	// Create temporary HTML file
	tempDir := "/tmp/pdf-conversion"
	os.MkdirAll(tempDir, 0755)

	timestamp := time.Now().Format("20060102_150405")
	htmlFile := filepath.Join(tempDir, fmt.Sprintf("input_%s.html", timestamp))
	pdfFile := filepath.Join(tempDir, fmt.Sprintf("output_%s.pdf", timestamp))

	// Write HTML content to file
	if err := os.WriteFile(htmlFile, []byte(htmlContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write HTML file: %w", err)
	}
	defer os.Remove(htmlFile)

	// Use Puppeteer/Chrome headless to generate PDF
	// This is a simplified approach - in production you'd want to use a proper headless browser
	// For now, we'll use a simple HTML to PDF conversion
	pdfData, err := convertHTMLToPDF(htmlFile, pdfFile, title)
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTML to PDF: %w", err)
	}

	// Clean up PDF file
	defer os.Remove(pdfFile)

	return pdfData, nil
}

func convertHTMLToPDF(htmlFile, pdfFile, title string) ([]byte, error) {
	// Use a simpler approach with chromium directly
	// This provides good PDF quality with UTF-8 support
	cmd := exec.Command("chromium-browser",
		"--headless",
		"--no-sandbox",
		"--disable-setuid-sandbox",
		"--disable-dev-shm-usage",
		"--disable-gpu",
		"--print-to-pdf="+pdfFile,
		"--print-to-pdf-no-header",
		"--run-all-compositor-stages-before-draw",
		"--virtual-time-budget=5000",
		"file://"+htmlFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("chromium failed: %w, output: %s", err, string(output))
	}

	// Read the generated PDF
	return os.ReadFile(pdfFile)
}
