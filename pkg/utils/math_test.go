package utils

import (
	"testing"
)

func TestGCD(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"basic", 48, 18, 6},
		{"one zero", 0, 5, 5},
		{"both zero", 0, 0, 0},
		{"coprime", 17, 31, 1},
		{"multiple", 100, 25, 25},
		{"negative", -48, 18, 6},
		{"both negative", -48, -18, -6},
		{"large numbers", 123456, 789012, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GCD(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("GCD(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestGCDOfSlice(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"empty slice", []int{}, 0},
		{"single element", []int{42}, 42},
		{"multiple elements", []int{48, 18, 30}, 6},
		{"coprime elements", []int{17, 31, 23}, 1},
		{"with zeros", []int{0, 5, 10}, 5},
		{"all zeros", []int{0, 0, 0}, 0},
		{"package sizes", []int{250, 500, 1000, 2000, 5000}, 250},
		{"divisible by common factor", []int{6, 12, 18, 24}, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GCDOfSlice(tt.numbers)
			if result != tt.expected {
				t.Errorf("GCDOfSlice(%v) = %d, want %d", tt.numbers, result, tt.expected)
			}
		})
	}
}

func TestValidatePackSizes(t *testing.T) {
	tests := []struct {
		name      string
		sizes     []int
		expectErr bool
		errMsg    string
	}{
		{"empty slice", []int{}, true, "pack sizes cannot be empty"},
		{"valid sizes", []int{250, 500, 1000}, false, ""},
		{"single valid size", []int{100}, false, ""},
		{"contains zero", []int{250, 0, 500}, true, "pack sizes must be positive integers"},
		{"contains negative", []int{250, -100, 500}, true, "pack sizes must be positive integers"},
		{"all negative", []int{-1, -2, -3}, true, "pack sizes must be positive integers"},
		{"contains one", []int{1, 250, 500}, false, ""},
		{"large numbers", []int{10000, 20000, 50000}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePackSizes(tt.sizes)
			if tt.expectErr {
				if err == nil {
					t.Errorf("ValidatePackSizes(%v) expected error, got nil", tt.sizes)
				} else if err.Error() != tt.errMsg {
					t.Errorf("ValidatePackSizes(%v) = %v, want %v", tt.sizes, err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePackSizes(%v) = %v, want nil", tt.sizes, err)
				}
			}
		})
	}
}
