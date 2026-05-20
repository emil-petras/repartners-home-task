package services

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

// TemplateData represents the data structure for template rendering
type TemplateData struct {
	Title          string
	PackSizes      []PackSizeDisplay
	SizesDisplay   string
	ErrorMessage   string
	SuccessMessage string
	ItemsValue     string
	FormAction     string
}

// PackSizeDisplay represents pack size for template display
type PackSizeDisplay struct {
	Size int
}

// TemplateService handles template rendering with proper separation of concerns
type TemplateService struct {
	templates *template.Template
}

// NewTemplateService creates a new template service with embedded templates
func NewTemplateService(webFS fs.FS) (*TemplateService, error) {
	// Parse templates from embedded filesystem
	templates, err := template.ParseFS(webFS, "*.html")
	if err != nil {
		return nil, err
	}

	return &TemplateService{
		templates: templates,
	}, nil
}

// RenderTemplate renders a template with the provided data
func (ts *TemplateService) RenderTemplate(w http.ResponseWriter, templateName string, data TemplateData) error {
	// Set content type header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute template with auto-escaping
	return ts.templates.ExecuteTemplate(w, templateName, data)
}

// RenderIndex renders the main index page with pack size data
func (ts *TemplateService) RenderIndex(w http.ResponseWriter, packSizes []int, errorMsg, successMsg, itemsValue string) error {
	// Convert pack sizes to display format
	var packSizeDisplays []PackSizeDisplay
	var sizesDisplay string

	if len(packSizes) == 0 {
		sizesDisplay = "No pack sizes set yet."
	} else {
		packSizeDisplays = make([]PackSizeDisplay, len(packSizes))
		sizes := make([]string, len(packSizes))
		for i, size := range packSizes {
			packSizeDisplays[i] = PackSizeDisplay{Size: size}
			sizes[i] = string(rune(size)) // This will be fixed in template
		}
		sizesDisplay = joinIntsToStrings(packSizes)
	}

	data := TemplateData{
		Title:          "Order Packs Calculator",
		PackSizes:      packSizeDisplays,
		SizesDisplay:   sizesDisplay,
		ErrorMessage:   errorMsg,
		SuccessMessage: successMsg,
		ItemsValue:     itemsValue,
		FormAction:     "/",
	}

	return ts.RenderTemplate(w, "simple.html", data)
}

// Helper function to convert int slice to comma-separated string
func joinIntsToStrings(ints []int) string {
	if len(ints) == 0 {
		return ""
	}

	strs := make([]string, len(ints))
	for i, num := range ints {
		strs[i] = fmt.Sprintf("%d", num)
	}
	return strings.Join(strs, ", ")
}
