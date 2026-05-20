package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"repartners-home-task/internal/database"
	"repartners-home-task/internal/models"
	"repartners-home-task/internal/services"
	web "repartners-home-task/web"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestHandler(t *testing.T) (*PackSizeHandler, func()) {
	// create temporary database for testing
	dbFile := fmt.Sprintf("test_db_%d.db", time.Now().UnixNano())
	db, err := database.NewDatabase(dbFile)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// create test services
	packSizeService := services.NewPackSizeService(db)
	packagingService := services.NewPackagingService(packSizeService)
	handler, err := NewPackSizeHandler(packSizeService, packagingService, web.FS)
	if err != nil {
		t.Fatalf("Failed to create test handler: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(dbFile)
	}

	return handler, cleanup
}

func TestPackSizeHandler_GetAllPackSizes(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// create test request
	req := httptest.NewRequest("GET", "/api/pack-sizes", nil)
	w := httptest.NewRecorder()

	handler.GetAllPackSizes(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// parse response
	var packSizes []models.PackSize
	err := json.NewDecoder(resp.Body).Decode(&packSizes)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// should be empty initially
	if len(packSizes) != 0 {
		t.Errorf("Expected 0 pack sizes, got %d", len(packSizes))
	}
}

func TestPackSizeHandler_ReplacePackSizes(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	tests := []struct {
		name           string
		requestBody    models.PackSizeRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid pack sizes",
			requestBody: models.PackSizeRequest{
				Sizes: []int{250, 500, 1000},
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "empty pack sizes",
			requestBody: models.PackSizeRequest{
				Sizes: []int{},
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "invalid JSON",
			requestBody: models.PackSizeRequest{
				Sizes: []int{0, -100},
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/api/pack-sizes", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.ReplacePackSizes(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.expectError {
				var errorResp map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&errorResp)
				if err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}

				if _, exists := errorResp["error"]; !exists {
					t.Error("Expected error field in response")
				}

				if _, exists := errorResp["code"]; !exists {
					t.Error("Expected code field in response")
				}
			}
		})
	}
}

func TestPackSizeHandler_CalculatePackaging(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// first setup some pack sizes
	packSizeReq := models.PackSizeRequest{Sizes: []int{250, 500, 1000}}
	body, _ := json.Marshal(packSizeReq)
	req := httptest.NewRequest("PUT", "/api/pack-sizes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ReplacePackSizes(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Failed to setup pack sizes: %d", w.Result().StatusCode)
	}

	tests := []struct {
		name           string
		requestBody    models.PackagingRequest
		expectedStatus int
		expectError    bool
	}{
		{
			name: "valid packaging request",
			requestBody: models.PackagingRequest{
				Items: 750,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "zero items",
			requestBody: models.PackagingRequest{
				Items: 0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "negative items",
			requestBody: models.PackagingRequest{
				Items: -100,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "large number of items",
			requestBody: models.PackagingRequest{
				Items: 5000,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/package", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			handler.CalculatePackaging(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.expectError {
				var errorResp map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&errorResp)
				if err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}

				if _, exists := errorResp["error"]; !exists {
					t.Error("Expected error field in response")
				}

				if _, exists := errorResp["code"]; !exists {
					t.Error("Expected code field in response")
				}
			} else {
				var result models.PackagingResponse
				err := json.NewDecoder(resp.Body).Decode(&result)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				// basic validation
				if result.Items != tt.requestBody.Items {
					t.Errorf("Items = %d, want %d", result.Items, tt.requestBody.Items)
				}

				if result.TotalItems < tt.requestBody.Items {
					t.Errorf("TotalItems = %d, should be >= %d", result.TotalItems, tt.requestBody.Items)
				}
			}
		})
	}
}

func TestPackSizeHandler_ServeWebPage(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/?success=Test%20success&items=500", nil)
	w := httptest.NewRecorder()

	handler.ServeWebPage(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("Response should contain HTML content")
	}
}

func TestPackSizeHandler_HandlePackSizeForm(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	tests := []struct {
		name        string
		formData    string
		expectError bool
	}{
		{
			name:        "valid form data",
			formData:    "sizes=250,500,1000",
			expectError: false,
		},
		{
			name:        "empty form data",
			formData:    "sizes=",
			expectError: true,
		},
		{
			name:        "invalid numbers",
			formData:    "sizes=250,abc,500",
			expectError: true,
		},
		{
			name:        "negative numbers",
			formData:    "sizes=250,-500,1000",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/pack-sizes", strings.NewReader(tt.formData))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.HandlePackSizeForm(w, req)

			resp := w.Result()
			// Should redirect (302) regardless of success or error
			if resp.StatusCode != http.StatusSeeOther {
				t.Errorf("Expected redirect status 303, got %d", resp.StatusCode)
			}

			// Check redirect location
			location := resp.Header.Get("Location")
			if location == "" {
				t.Error("Expected Location header in redirect response")
			}
		})
	}
}

func TestPackSizeHandler_HandlePackagingForm(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// Setup some pack sizes first
	packSizeReq := models.PackSizeRequest{Sizes: []int{250, 500, 1000}}
	body, _ := json.Marshal(packSizeReq)
	req := httptest.NewRequest("PUT", "/api/pack-sizes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ReplacePackSizes(w, req)

	tests := []struct {
		name        string
		formData    string
		expectError bool
	}{
		{
			name:        "valid form data",
			formData:    "items=500",
			expectError: false,
		},
		{
			name:        "empty form data",
			formData:    "items=",
			expectError: true,
		},
		{
			name:        "invalid numbers",
			formData:    "items=abc",
			expectError: true,
		},
		{
			name:        "negative numbers",
			formData:    "items=-100",
			expectError: true,
		},
		{
			name:        "zero items",
			formData:    "items=0",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/package", strings.NewReader(tt.formData))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()

			handler.HandlePackagingForm(w, req)

			resp := w.Result()
			// Should redirect (302) regardless of success or error
			if resp.StatusCode != http.StatusSeeOther {
				t.Errorf("Expected redirect status 303, got %d", resp.StatusCode)
			}

			// Check redirect location
			location := resp.Header.Get("Location")
			if location == "" {
				t.Error("Expected Location header in redirect response")
			}
		})
	}
}
