package services

import (
	"fmt"
	"math"
	"repartners-home-task/internal/models"
	"repartners-home-task/pkg/utils"
	"sort"
)

type PackagingService struct {
	packSizeService *PackSizeService
}

func NewPackagingService(packSizeService *PackSizeService) *PackagingService {
	return &PackagingService{packSizeService: packSizeService}
}

// DPState represents the state in our dynamic programming solution
type DPState struct {
	packages int // number of packages used
	lastUsed int // size of last pack used
}

// CalculateOptimalPackaging calculates the optimal packaging for given items
func (s *PackagingService) CalculateOptimalPackaging(items int) (*models.PackagingResponse, error) {
	// validate input
	if items <= 0 {
		return nil, fmt.Errorf("items must be positive")
	}

	// get available pack sizes from database
	packSizes, err := s.packSizeService.GetAllPackSizes()
	if err != nil {
		return nil, fmt.Errorf("failed to get pack sizes: %w", err)
	}

	// check if any pack sizes are configured
	if len(packSizes) == 0 {
		return nil, fmt.Errorf("no pack sizes configured")
	}

	// extract pack sizes into int slice
	sizes := make([]int, len(packSizes))
	for i, ps := range packSizes {
		sizes[i] = ps.Size
	}

	// validate pack sizes
	if err := utils.ValidatePackSizes(sizes); err != nil {
		return nil, err
	}

	// sort pack sizes for consistent processing
	sort.Ints(sizes)

	// find gcd of all pack sizes for optimization
	gcd := utils.GCDOfSlice(sizes)
	normalizedSizes := make([]int, len(sizes))
	for i, size := range sizes {
		normalizedSizes[i] = size / gcd
	}

	// normalize items count
	normalizedItems := items / gcd
	if items%gcd != 0 {
		normalizedItems++
	}

	// calculate maximum total items we need to consider
	maxSize := normalizedSizes[len(normalizedSizes)-1]
	maxTotal := normalizedItems + maxSize - 1

	// initialize dynamic programming table
	dp := make([]DPState, maxTotal+1)
	for i := range dp {
		dp[i] = DPState{packages: math.MaxInt32, lastUsed: -1}
	}
	dp[0] = DPState{packages: 0, lastUsed: 0}

	// fill dp table using dynamic programming
	for _, packSize := range normalizedSizes {
		for total := packSize; total <= maxTotal; total++ {
			prev := total - packSize
			if dp[prev].packages != math.MaxInt32 {
				if dp[prev].packages+1 < dp[total].packages {
					dp[total] = DPState{
						packages: dp[prev].packages + 1,
						lastUsed: packSize,
					}
				}
			}
		}
	}

	// find the best total (minimum overage)
	bestTotal := -1
	for total := normalizedItems; total <= maxTotal; total++ {
		if dp[total].packages != math.MaxInt32 {
			bestTotal = total
			break
		}
	}

	// check if solution was found
	if bestTotal == -1 {
		return nil, fmt.Errorf("no valid packaging combination found")
	}

	// backtrack to find package counts
	packageCounts := make(map[int]int)
	current := bestTotal

	for current > 0 {
		packSize := dp[current].lastUsed
		originalSize := packSize * gcd
		packageCounts[originalSize]++
		current -= packSize
	}

	// calculate actual total items
	actualTotal := bestTotal * gcd

	// create package details for response
	details := make([]models.PackageDetail, 0)
	for size, count := range packageCounts {
		details = append(details, models.PackageDetail{
			Size:  size,
			Count: count,
		})
	}

	// sort details by pack size (descending)
	sort.Slice(details, func(i, j int) bool {
		return details[i].Size > details[j].Size
	})

	// return the packaging response
	return &models.PackagingResponse{
		Items:      items,
		TotalItems: actualTotal,
		Packages:   packageCounts,
		Details:    details,
	}, nil
}
