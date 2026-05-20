package handlers

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"repartners-home-task/internal/models"
	"repartners-home-task/internal/services"

	"github.com/go-chi/chi/v5"
)

// PackSizeHandler manages HTTP requests for pack sizes and packaging calculations.
type PackSizeHandler struct {
	service          *services.PackSizeService
	packagingService *services.PackagingService
	templateService  *services.TemplateService
}

// NewPackSizeHandler initializes a handler with the required services.
func NewPackSizeHandler(service *services.PackSizeService, packagingService *services.PackagingService, webFS fs.FS) (*PackSizeHandler, error) {
	templateService, err := services.NewTemplateService(webFS)
	if err != nil {
		return nil, err
	}

	return &PackSizeHandler{
		service:          service,
		packagingService: packagingService,
		templateService:  templateService,
	}, nil
}

// RegisterRoutes configures all application routes.
func (h *PackSizeHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.ServeWebPage)
	r.Post("/pack-sizes", h.HandlePackSizeForm)
	r.Post("/package", h.HandlePackagingForm)

	r.Route("/api/pack-sizes", func(r chi.Router) {
		r.Get("/", h.GetAllPackSizes)
		r.Put("/", h.ReplacePackSizes)
	})

	r.Post("/api/package", h.CalculatePackaging)
}

func (h *PackSizeHandler) GetAllPackSizes(w http.ResponseWriter, r *http.Request) {
	packSizes, err := h.service.GetAllPackSizes()
	if err != nil {
		sendErrorResponse(w, r, http.StatusInternalServerError, "failed to retrieve pack sizes", err.Error())
		return
	}

	sendSuccessResponse(w, r, packSizes)
}

func (h *PackSizeHandler) ReplacePackSizes(w http.ResponseWriter, r *http.Request) {
	var req models.PackSizeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, r, http.StatusBadRequest, "invalid json format", err.Error())
		return
	}

	if len(req.Sizes) == 0 {
		sendErrorResponse(w, r, http.StatusBadRequest, "pack sizes array cannot be empty")
		return
	}

	for _, size := range req.Sizes {
		if size <= 0 {
			sendErrorResponse(w, r, http.StatusBadRequest, "pack sizes must be positive integers")
			return
		}
	}

	if err := h.service.ReplacePackSizes(req.Sizes); err != nil {
		sendErrorResponse(w, r, http.StatusInternalServerError, "failed to replace pack sizes", err.Error())
		return
	}

	sendSuccessResponse(w, r, map[string]interface{}{
		"message": "pack sizes replaced successfully",
		"count":   len(req.Sizes),
	})
}

func (h *PackSizeHandler) CalculatePackaging(w http.ResponseWriter, r *http.Request) {
	var req models.PackagingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, r, http.StatusBadRequest, "invalid json format", err.Error())
		return
	}

	if req.Items <= 0 {
		sendErrorResponse(w, r, http.StatusBadRequest, "items must be a positive integer")
		return
	}

	result, err := h.packagingService.CalculateOptimalPackaging(req.Items)
	if err != nil {
		sendErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	sendSuccessResponse(w, r, result)
}

// ServeWebPage renders the main HTML interface.
func (h *PackSizeHandler) ServeWebPage(w http.ResponseWriter, r *http.Request) {
	packSizes, err := h.service.GetAllPackSizes()
	if err != nil {
		http.Error(w, "Failed to load pack sizes", http.StatusInternalServerError)
		return
	}

	sizes := make([]int, len(packSizes))
	for i, ps := range packSizes {
		sizes[i] = ps.Size
	}

	errorMsg := r.URL.Query().Get("error")
	successMsg := r.URL.Query().Get("success")
	itemsValue := r.URL.Query().Get("items")

	if err = h.templateService.RenderIndex(w, sizes, errorMsg, successMsg, itemsValue); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// HandlePackSizeForm processes pack size form submissions.
func (h *PackSizeHandler) HandlePackSizeForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		redirectWithError(w, r, "Failed to parse form data")
		return
	}

	sizesStr := r.FormValue("sizes")
	if sizesStr == "" {
		redirectWithError(w, r, "Pack sizes cannot be empty")
		return
	}

	sizes, err := parseSizes(sizesStr)
	if err != nil {
		redirectWithError(w, r, fmt.Sprintf("Invalid pack sizes: %v", err))
		return
	}

	if len(sizes) == 0 {
		redirectWithError(w, r, "At least one pack size is required")
		return
	}

	for _, size := range sizes {
		if size <= 0 {
			redirectWithError(w, r, "Pack sizes must be positive")
			return
		}
	}

	if err := h.service.ReplacePackSizes(sizes); err != nil {
		redirectWithError(w, r, "Failed to save pack sizes")
		return
	}

	redirectWithMessage(w, r, "Pack sizes updated successfully")
}

// HandlePackagingForm processes packaging calculation requests from the web form.
func (h *PackSizeHandler) HandlePackagingForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		redirectWithError(w, r, "Failed to parse form data")
		return
	}

	itemsStr := r.FormValue("items")
	if itemsStr == "" {
		redirectWithError(w, r, "Items cannot be empty")
		return
	}

	var items int
	if _, err := fmt.Sscanf(itemsStr, "%d", &items); err != nil || items <= 0 {
		redirectWithError(w, r, "Items must be a positive number")
		return
	}

	result, err := h.packagingService.CalculateOptimalPackaging(items)
	if err != nil {
		redirectWithError(w, r, fmt.Sprintf("Packaging calculation failed: %v", err))
		return
	}

	var parts []string
	for size, count := range result.Packages {
		parts = append(parts, fmt.Sprintf("%d×%d", count, size))
	}

	msg := fmt.Sprintf("Packaging calculated for %d items. Total shipped: %d. Packages: %s",
		items, result.TotalItems, strings.Join(parts, ", "))

	http.Redirect(w, r, fmt.Sprintf("/?success=%s&items=%d", msg, items), http.StatusSeeOther)
}
