package handlers

import (
	"fmt"
	"strings"
)

// parseSizes parses a comma-separated string of sizes into a slice of integers
func parseSizes(sizesStr string) ([]int, error) {
	// proper comma-separated parser
	var sizes []int
	for _, part := range strings.Split(sizesStr, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		var size int
		if _, err := fmt.Sscanf(part, "%d", &size); err != nil {
			return nil, fmt.Errorf("invalid number: %s", part)
		}
		sizes = append(sizes, size)
	}
	return sizes, nil
}
