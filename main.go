package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"pdffear/helper"
	"strings"
)

var TMP_PDF_PATH = os.Getenv("TMP_PDF_PATH")
var TMP_DOCX_PATH = os.Getenv("TMP_DOCX_PATH")
var LIBREOFFICE_PROFILES = os.Getenv("LIBREOFFICE_PROFILES") + "/" + helper.GenerateRandomString(10)

func downloadFile(filename string, w http.ResponseWriter, r *http.Request) error {
	filename = strings.TrimSuffix(filename, ".docx")
	convertedFile := filepath.Join(TMP_PDF_PATH, filename + ".pdf") 
	
	file, err := os.Open(convertedFile)
	if err != nil {
		http.Error(w, "Error Opening the File", http.StatusInternalServerError)
		return err
	}
	defer file.Close()
	
	info, err := file.Stat()
	if err != nil {
		// Handle error
		log.Fatal(err)
		return err
	}
	
	w.Header().Set("Content-Disposition", "attachment; filename=" + filename)
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeContent(w, r, filename, info.ModTime(), file)
	return nil
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return 
	}

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := handler.Filename
	docxFilename := handler.Filename
	destinationPath := filepath.Join(TMP_DOCX_PATH, docxFilename)

	dst, err := os.Create(".." + destinationPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error Uploading the File %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error Copying the File", http.StatusInternalServerError)
		return
	}


	err = convertToPdf(docxFilename, w)
	if err != nil {
		http.Error(w, "Error Converting the File", http.StatusInternalServerError)
		return
	}
	downloadFile(filename, w, r)
}

func convertToPdf(docxFilename string, w http.ResponseWriter) error {
	cmd := "libreoffice"
    args := []string{
        fmt.Sprintf("-env:UserInstallation=file://%s", LIBREOFFICE_PROFILES),
        "--headless",
        "--convert-to", "pdf:writer_pdf_Export",
  		"--outdir", TMP_PDF_PATH,
        TMP_DOCX_PATH + docxFilename,
    }

	output, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		w.Write([]byte(output))
		http.Error(w, fmt.Sprintf("Error Running the Command %s", err.Error()), http.StatusInternalServerError)
		return err
	} 

	return nil
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Server is up and running"))
}

func main() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/health", healthCheck)
	http.ListenAndServe(":8080", nil)
}

