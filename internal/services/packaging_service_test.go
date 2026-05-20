package services

import (
	"os"
	"testing"

	"repartners-home-task/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestPackagingService(t *testing.T) (*PackagingService, func()) {
	// create temporary database
	dbFile := "test_packaging_" + t.Name() + ".db"

	db, err := database.NewDatabase(dbFile)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// setup initial pack sizes
	testSizes := []int{250, 500, 1000, 2000, 5000}
	err = db.ReplacePackSizes(testSizes)
	if err != nil {
		t.Fatalf("Failed to setup test pack sizes: %v", err)
	}

	packSizeService := NewPackSizeService(db)
	packagingService := NewPackagingService(packSizeService)

	cleanup := func() {
		db.Close()
		os.Remove(dbFile)
	}

	return packagingService, cleanup
}

func TestPackagingService_CalculateOptimalPackaging(t *testing.T) {
	packagingService, cleanup := setupTestPackagingService(t)
	defer cleanup()

	tests := []struct {
		name          string
		items         int
		expectedTotal int
		expectedPacks map[int]int
		expectErr     bool
		errMsg        string
	}{
		{
			name:          "exact match",
			items:         500,
			expectedTotal: 500,
			expectedPacks: map[int]int{500: 1},
		},
		{
			name:          "needs rounding up",
			items:         750,
			expectedTotal: 750,
			expectedPacks: map[int]int{250: 1, 500: 1},
		},
		{
			name:          "small order",
			items:         1,
			expectedTotal: 250,
			expectedPacks: map[int]int{250: 1},
		},
		{
			name:          "large order",
			items:         1200,
			expectedTotal: 1250,
			expectedPacks: map[int]int{250: 1, 1000: 1},
		},
		{
			name:          "very large order",
			items:         5300,
			expectedTotal: 5500,
			expectedPacks: map[int]int{500: 1, 5000: 1},
		},
		{
			name:          "multiple small packs",
			items:         3000,
			expectedTotal: 3000,
			expectedPacks: map[int]int{2000: 1, 1000: 1},
		},
		{
			name:          "optimized for fewer packs",
			items:         2600,
			expectedTotal: 2750,
			expectedPacks: map[int]int{2000: 1, 500: 1, 250: 1},
		},
		{
			name:      "zero items",
			items:     0,
			expectErr: true,
			errMsg:    "items must be positive",
		},
		{
			name:      "negative items",
			items:     -100,
			expectErr: true,
			errMsg:    "items must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := packagingService.CalculateOptimalPackaging(tt.items)

			if tt.expectErr {
				if err == nil {
					t.Errorf("CalculateOptimalPackaging(%d) expected error, got nil", tt.items)
				} else if err.Error() != tt.errMsg {
					t.Errorf("CalculateOptimalPackaging(%d) = %v, want %v", tt.items, err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("CalculateOptimalPackaging(%d) = %v, want nil", tt.items, err)
				return
			}

			// check total items
			if result.TotalItems != tt.expectedTotal {
				t.Errorf("TotalItems = %d, want %d", result.TotalItems, tt.expectedTotal)
			}

			// check requested items
			if result.Items != tt.items {
				t.Errorf("Items = %d, want %d", result.Items, tt.items)
			}

			// check package counts
			if len(result.Packages) != len(tt.expectedPacks) {
				t.Errorf("Packages length = %d, want %d", len(result.Packages), len(tt.expectedPacks))
			}

			for size, expectedCount := range tt.expectedPacks {
				if count, exists := result.Packages[size]; !exists {
					t.Errorf("Package size %d not found in result", size)
				} else if count != expectedCount {
					t.Errorf("Package size %d count = %d, want %d", size, count, expectedCount)
				}
			}

			// verify details match packages
			if len(result.Details) != len(result.Packages) {
				t.Errorf("Details length = %d, packages length = %d", len(result.Details), len(result.Packages))
			}

			// verify total items calculation
			calculatedTotal := 0
			for size, count := range result.Packages {
				calculatedTotal += size * count
			}
			if calculatedTotal != result.TotalItems {
				t.Errorf("Calculated total = %d, result.TotalItems = %d", calculatedTotal, result.TotalItems)
			}
		})
	}
}

func TestPackagingService_WithDifferentPackSizes(t *testing.T) {
	// test with different pack size configurations
	testCases := []struct {
		name      string
		packSizes []int
		testItems int
		expectErr bool
	}{
		{"small packs", []int{1, 2, 5}, 7, false},
		{"large packs", []int{1000, 2000, 5000}, 3000, false},
		{"single pack", []int{100}, 250, false},
		{"coprime packs", []int{7, 11, 13}, 20, false},
		{"divisible packs", []int{3, 6, 9, 12}, 10, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// create temporary database with specific pack sizes
			dbFile := "test_custom_" + tc.name + ".db"
			db, err := database.NewDatabase(dbFile)
			if err != nil {
				t.Fatalf("Failed to create test database: %v", err)
			}
			defer func() {
				db.Close()
				os.Remove(dbFile)
			}()

			// setup pack sizes
			err = db.ReplacePackSizes(tc.packSizes)
			if err != nil {
				t.Fatalf("Failed to setup pack sizes: %v", err)
			}

			packSizeService := NewPackSizeService(db)
			packagingService := NewPackagingService(packSizeService)

			result, err := packagingService.CalculateOptimalPackaging(tc.testItems)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected error for pack sizes %v and items %d", tc.packSizes, tc.testItems)
				}
				return
			}

			if err != nil {
				t.Errorf("CalculateOptimalPackaging failed: %v", err)
				return
			}

			// basic validation
			if result.Items != tc.testItems {
				t.Errorf("Items = %d, want %d", result.Items, tc.testItems)
			}

			if result.TotalItems < tc.testItems {
				t.Errorf("TotalItems = %d, should be >= %d", result.TotalItems, tc.testItems)
			}

			// verify we can reconstruct the total
			calculatedTotal := 0
			for size, count := range result.Packages {
				calculatedTotal += size * count
			}
			if calculatedTotal != result.TotalItems {
				t.Errorf("Total calculation mismatch: %d vs %d", calculatedTotal, result.TotalItems)
			}
		})
	}
}

func TestPackagingService_NoPackSizes(t *testing.T) {
	// test behavior when no pack sizes are configured
	dbFile := "test_empty_" + t.Name() + ".db"
	db, err := database.NewDatabase(dbFile)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer func() {
		db.Close()
		os.Remove(dbFile)
	}()

	packSizeService := NewPackSizeService(db)
	packagingService := NewPackagingService(packSizeService)

	_, err = packagingService.CalculateOptimalPackaging(100)
	if err == nil {
		t.Error("Expected error when no pack sizes configured")
	}

	if err.Error() != "no pack sizes configured" {
		t.Errorf("Expected 'no pack sizes configured', got %v", err.Error())
	}
}

func TestPackagingService_OptimizationPriority(t *testing.T) {
	packagingService, cleanup := setupTestPackagingService(t)
	defer cleanup()

	// test that the algorithm prioritizes minimizing total items over minimizing packages
	// for 2600 items with packs [250, 500, 1000, 2000, 5000]:
	// option 1: 2000 + 500 + 250 = 2750 (3 packs)
	// option 2: 1000 + 1000 + 1000 = 3000 (3 packs)
	// option 3: 5000 = 5000 (1 pack)
	// the algorithm should choose option 1 (minimizes total items)

	result, err := packagingService.CalculateOptimalPackaging(2600)
	if err != nil {
		t.Fatalf("CalculateOptimalPackaging failed: %v", err)
	}

	// should choose 2750 total items (minimum overage)
	if result.TotalItems != 2750 {
		t.Errorf("Expected total items 2750, got %d", result.TotalItems)
	}

	// should use 3 packages (2000, 500, 250)
	totalPackages := 0
	for _, count := range result.Packages {
		totalPackages += count
	}

	if totalPackages != 3 {
		t.Errorf("Expected 3 packages, got %d", totalPackages)
	}

	// verify the specific pack sizes
	expectedPacks := map[int]int{2000: 1, 500: 1, 250: 1}
	for size, expectedCount := range expectedPacks {
		if count, exists := result.Packages[size]; !exists {
			t.Errorf("Expected pack size %d not found", size)
		} else if count != expectedCount {
			t.Errorf("Pack size %d count = %d, want %d", size, count, expectedCount)
		}
	}
}
