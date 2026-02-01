package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"pdffear/converter"
	"pdffear/helper"
	"pdffear/storage"
)

// Config holds the application configuration.
type Config struct {
	TmpPdfPath         string
	TmpDocxPath        string
	LibreOfficeProfile string
}

// App holds the application dependencies.
type App struct {
	Config    *Config
	Converter converter.Converter
	Storage   storage.Storage
}

func main() {
	// Initialize Configuration
	config := &Config{
		TmpPdfPath:         getEnv("TMP_PDF_PATH", "/tmp/generated_pdfs"),
		TmpDocxPath:        getEnv("TMP_DOCX_PATH", "/tmp/uploaded_docx"),
		LibreOfficeProfile: getEnv("LIBREOFFICE_PROFILES", "/tmp/libreoffice_profiles") + helper.GenerateRandomString(10),
	}

	// Initialize Directories
	ensureDir(config.TmpPdfPath)
	ensureDir(config.TmpDocxPath)
	ensureDir(config.LibreOfficeProfile)

	// Initialize Services
	app := &App{
		Config:    config,
		Converter: converter.NewLibreOfficeConverter("libreoffice", config.LibreOfficeProfile),
		Storage:   storage.NewLocalStorage(),
	}

	// Setup Routes
	http.HandleFunc("/upload", app.uploadHandler)
	http.HandleFunc("/health", app.healthHandler)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// uploadHandler handles file uploads and conversion.
func (app *App) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse Form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error Parsing Form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 1. Save File
	savedPath, err := app.Storage.Save(file, handler, app.Config.TmpDocxPath)
	if err != nil {
		log.Printf("Error saving file: %v", err)
		http.Error(w, fmt.Sprintf("Error Saving File: %v", err), http.StatusInternalServerError)
		return
	}

	// 2. Convert File
	pdfPath, err := app.Converter.Convert(savedPath, app.Config.TmpPdfPath)
	if err != nil {
		log.Printf("Error converting file: %v", err)
		http.Error(w, fmt.Sprintf("Error Converting File: %v", err), http.StatusInternalServerError)
		return
	}

	// 3. Download File
	app.downloadFile(w, r, pdfPath)
}

// downloadFile serves the converted PDF file.
func (app *App) downloadFile(w http.ResponseWriter, r *http.Request, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening PDF: %v", err)
		http.Error(w, "Error Opening Converted File", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Printf("Error stating PDF: %v", err)
		http.Error(w, "Error Stating File", http.StatusInternalServerError)
		return
	}

	filename := filepath.Base(filePath)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeContent(w, r, filename, info.ModTime(), file)
}

// healthHandler checks if the server is running.
func (app *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Server is up and running"))
}

// Helper functions

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ensureDir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("Warning: ensureDir failed for %s: %v", path, err)
	}
}
