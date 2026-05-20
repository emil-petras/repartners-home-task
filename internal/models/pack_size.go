package models

import (
	"time"
)

// PackSize represents a pack size entity in the database
type PackSize struct {
	ID        int       `json:"id" db:"id"`                 // unique identifier
	Size      int       `json:"size" db:"size"`             // pack size value
	CreatedAt time.Time `json:"created_at" db:"created_at"` // creation timestamp
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // last update timestamp
}

// PackSizeRequest represents a request to replace pack sizes
type PackSizeRequest struct {
	Sizes []int `json:"sizes" validate:"required"` // array of pack sizes
}

// PackSizeResponse represents a pack size in api responses
type PackSizeResponse struct {
	ID        int       `json:"id"`         // unique identifier
	Size      int       `json:"size"`       // pack size value
	CreatedAt time.Time `json:"created_at"` // creation timestamp
	UpdatedAt time.Time `json:"updated_at"` // last update timestamp
}

// PackagingRequest represents a request to calculate packaging
type PackagingRequest struct {
	Items int `json:"items" validate:"required,min=1"` // number of items to package
}

// PackagingResponse represents the result of packaging calculation
type PackagingResponse struct {
	Items      int             `json:"items"`       // original number of items
	TotalItems int             `json:"total_items"` // total items shipped (may be more than requested)
	Packages   map[int]int     `json:"packages"`    // map of pack size to quantity
	Details    []PackageDetail `json:"details"`     // detailed breakdown of packages
}

// PackageDetail represents a single package type in the result
type PackageDetail struct {
	Size  int `json:"size"`  // pack size
	Count int `json:"count"` // number of packages of this size
}
