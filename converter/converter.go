package converter

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Converter defines the interface for file conversion.
type Converter interface {
	Convert(inputPath string, outputDir string) (string, error)
}

// LibreOfficeConverter implements Converter using LibreOffice.
type LibreOfficeConverter struct {
	LibreOfficePath string
	ProfilePath     string
}

// NewLibreOfficeConverter creates a new LibreOfficeConverter.
func NewLibreOfficeConverter(libreOfficePath string, profilePath string) *LibreOfficeConverter {
	return &LibreOfficeConverter{
		LibreOfficePath: libreOfficePath,
		ProfilePath:     profilePath,
	}
}

// Convert converts a document to PDF.
func (c *LibreOfficeConverter) Convert(inputPath string, outputDir string) (string, error) {
	ext := strings.ToLower(filepath.Ext(inputPath))
	var exportFilter string

	switch ext {
	case ".docx", ".doc":
		exportFilter = "writer_pdf_Export"
	case ".xlsx", ".xls":
		exportFilter = "calc_pdf_Export"
	default:
		return "", fmt.Errorf("unsupported file extension: %s", ext)
	}

	args := []string{
		fmt.Sprintf("-env:UserInstallation=file://%s", c.ProfilePath),
		"--headless",
		"--convert-to", "pdf:" + exportFilter,
		"--outdir", outputDir,
		inputPath,
	}

	cmd := exec.Command(c.LibreOfficePath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("libreoffice conversion failed: %s, output: %s", err, string(output))
	}

	// Calculate output filename
	baseName := strings.TrimSuffix(filepath.Base(inputPath), ext)
	return filepath.Join(outputDir, baseName+".pdf"), nil
}
