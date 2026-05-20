package utils

import (
	"fmt"
)

// GCD calculates the greatest common divisor of two integers using euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// GCDOfSlice calculates the gcd of all integers in a slice
func GCDOfSlice(numbers []int) int {
	// handle empty slice
	if len(numbers) == 0 {
		return 0
	}

	// start with first number
	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		result = GCD(result, numbers[i])
		// early exit if gcd is 1 (can't get smaller)
		if result == 1 {
			break
		}
	}
	return result
}

// ValidatePackSizes validates that pack sizes are valid for packaging calculations
func ValidatePackSizes(sizes []int) error {
	// check if slice is empty
	if len(sizes) == 0 {
		return fmt.Errorf("pack sizes cannot be empty")
	}

	// check if all sizes are positive
	for _, size := range sizes {
		if size <= 0 {
			return fmt.Errorf("pack sizes must be positive integers")
		}
	}

	return nil
}
