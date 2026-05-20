package services

import (
	"fmt"
	"os"
	"testing"

	"repartners-home-task/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestPackSizeService(t *testing.T) (*PackSizeService, func()) {
	// create temporary database
	dbFile := "test_packsize_" + t.Name() + ".db"

	db, err := database.NewDatabase(dbFile)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	service := NewPackSizeService(db)

	cleanup := func() {
		db.Close()
		os.Remove(dbFile)
	}

	return service, cleanup
}

func TestNewPackSizeService(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	if service == nil {
		t.Fatal("NewPackSizeService returned nil")
	}

	if service.db == nil {
		t.Fatal("PackSizeService.db is nil")
	}
}

func TestPackSizeService_GetAllPackSizes_Empty(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	// test empty database
	sizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("GetAllPackSizes failed: %v", err)
	}

	if len(sizes) != 0 {
		t.Errorf("Expected 0 sizes, got %d", len(sizes))
	}
}

func TestPackSizeService_GetAllPackSizes_WithData(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	// setup test data
	testSizes := []int{250, 500, 1000, 2000}
	err := service.ReplacePackSizes(testSizes)
	if err != nil {
		t.Fatalf("ReplacePackSizes failed: %v", err)
	}

	// test retrieval
	sizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("GetAllPackSizes failed: %v", err)
	}

	if len(sizes) != len(testSizes) {
		t.Errorf("Expected %d sizes, got %d", len(testSizes), len(sizes))
	}

	// verify sizes match (order may vary)
	sizeMap := make(map[int]bool)
	for _, size := range testSizes {
		sizeMap[size] = true
	}

	for _, ps := range sizes {
		if !sizeMap[ps.Size] {
			t.Errorf("Unexpected size %d found", ps.Size)
		}
	}
}

func TestPackSizeService_ReplacePackSizes_New(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	testSizes := []int{100, 250, 500, 1000, 2000}

	// test replacement
	err := service.ReplacePackSizes(testSizes)
	if err != nil {
		t.Fatalf("ReplacePackSizes failed: %v", err)
	}

	// verify replacement
	sizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("GetAllPackSizes failed: %v", err)
	}

	if len(sizes) != len(testSizes) {
		t.Errorf("Expected %d sizes, got %d", len(testSizes), len(sizes))
	}
}

func TestPackSizeService_ReplacePackSizes_ReplaceExisting(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	// initial data
	initialSizes := []int{250, 500, 1000}
	err := service.ReplacePackSizes(initialSizes)
	if err != nil {
		t.Fatalf("Initial ReplacePackSizes failed: %v", err)
	}

	// replace with new data
	newSizes := []int{100, 200, 300, 400}
	err = service.ReplacePackSizes(newSizes)
	if err != nil {
		t.Fatalf("ReplacePackSizes failed: %v", err)
	}

	// verify replacement
	sizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("GetAllPackSizes failed: %v", err)
	}

	if len(sizes) != len(newSizes) {
		t.Errorf("Expected %d sizes, got %d", len(newSizes), len(sizes))
	}

	// verify old sizes are gone
	sizeMap := make(map[int]bool)
	for _, size := range newSizes {
		sizeMap[size] = true
	}

	for _, ps := range sizes {
		if !sizeMap[ps.Size] {
			t.Errorf("Old size %d still present", ps.Size)
		}
	}
}

func TestPackSizeService_ReplacePackSizes_Empty(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	// setup initial data
	initialSizes := []int{250, 500, 1000}
	err := service.ReplacePackSizes(initialSizes)
	if err != nil {
		t.Fatalf("Initial ReplacePackSizes failed: %v", err)
	}

	// replace with empty
	err = service.ReplacePackSizes([]int{})
	if err != nil {
		t.Fatalf("ReplacePackSizes with empty failed: %v", err)
	}

	// verify empty
	sizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("GetAllPackSizes failed: %v", err)
	}

	if len(sizes) != 0 {
		t.Errorf("Expected 0 sizes, got %d", len(sizes))
	}
}

func TestPackSizeService_ReplacePackSizes_Single(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	testSizes := []int{1000}

	// test single size
	err := service.ReplacePackSizes(testSizes)
	if err != nil {
		t.Fatalf("ReplacePackSizes failed: %v", err)
	}

	// verify single size
	sizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("GetAllPackSizes failed: %v", err)
	}

	if len(sizes) != 1 {
		t.Errorf("Expected 1 size, got %d", len(sizes))
	}

	if sizes[0].Size != 1000 {
		t.Errorf("Expected size 1000, got %d", sizes[0].Size)
	}
}

func TestPackSizeService_ReplacePackSizes_LargeNumbers(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	testSizes := []int{100000, 200000, 500000, 1000000}

	// test large numbers
	err := service.ReplacePackSizes(testSizes)
	if err != nil {
		t.Fatalf("ReplacePackSizes failed: %v", err)
	}

	// verify large numbers
	sizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("GetAllPackSizes failed: %v", err)
	}

	if len(sizes) != len(testSizes) {
		t.Errorf("Expected %d sizes, got %d", len(testSizes), len(sizes))
	}
}

func TestPackSizeService_ReplacePackSizes_Unsorted(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	// unsorted input
	testSizes := []int{1000, 250, 5000, 500, 2000}

	err := service.ReplacePackSizes(testSizes)
	if err != nil {
		t.Fatalf("ReplacePackSizes failed: %v", err)
	}

	// verify all sizes are present
	sizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("GetAllPackSizes failed: %v", err)
	}

	if len(sizes) != len(testSizes) {
		t.Errorf("Expected %d sizes, got %d", len(testSizes), len(sizes))
	}

	// verify all unique sizes are present
	sizeMap := make(map[int]bool)
	for _, size := range testSizes {
		sizeMap[size] = true
	}

	for _, ps := range sizes {
		if !sizeMap[ps.Size] {
			t.Errorf("Size %d not found in input", ps.Size)
		}
	}
}

func TestPackSizeService_ReplacePackSizes_Duplicates(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	// input with duplicates - service should handle this gracefully
	testSizes := []int{250, 500, 250, 1000, 500}

	err := service.ReplacePackSizes(testSizes)
	// duplicates should cause an error due to UNIQUE constraint
	if err == nil {
		t.Error("Expected error for duplicate pack sizes")
	}

	// verify database state after failed operation
	sizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("GetAllPackSizes failed: %v", err)
	}

	// should still be empty since the transaction failed
	if len(sizes) != 0 {
		t.Errorf("Expected 0 sizes after failed duplicate insert, got %d", len(sizes))
	}
}

func TestPackSizeService_Integration_GetAfterReplace(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	// multiple replace and get operations
	testCases := [][]int{
		{100, 200},
		{50, 100, 150, 200, 250},
		{1000},
		{},
		{75, 150, 300, 600},
	}

	for i, testSizes := range testCases {
		t.Run(fmt.Sprintf("iteration_%d", i), func(t *testing.T) {
			// replace
			err := service.ReplacePackSizes(testSizes)
			if err != nil {
				t.Fatalf("ReplacePackSizes failed: %v", err)
			}

			// get and verify
			sizes, err := service.GetAllPackSizes()
			if err != nil {
				t.Fatalf("GetAllPackSizes failed: %v", err)
			}

			if len(sizes) != len(testSizes) {
				t.Errorf("Expected %d sizes, got %d", len(testSizes), len(sizes))
			}

			// verify content
			sizeMap := make(map[int]bool)
			for _, size := range testSizes {
				sizeMap[size] = true
			}

			for _, ps := range sizes {
				if !sizeMap[ps.Size] {
					t.Errorf("Unexpected size %d", ps.Size)
				}
			}
		})
	}
}

func TestPackSizeService_ModelValidation(t *testing.T) {
	service, cleanup := setupTestPackSizeService(t)
	defer cleanup()

	testSizes := []int{250, 500, 1000}

	// replace pack sizes
	err := service.ReplacePackSizes(testSizes)
	if err != nil {
		t.Fatalf("ReplacePackSizes failed: %v", err)
	}

	// get pack sizes and validate model structure
	sizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("GetAllPackSizes failed: %v", err)
	}

	for i, ps := range sizes {
		if ps.Size <= 0 {
			t.Errorf("Size %d should be positive, got %d", i, ps.Size)
		}

		if ps.CreatedAt.IsZero() {
			t.Errorf("Size %d should have non-zero CreatedAt", i)
		}

		if ps.UpdatedAt.IsZero() {
			t.Errorf("Size %d should have non-zero UpdatedAt", i)
		}
	}
}
