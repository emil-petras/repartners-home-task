package services

import (
	"repartners-home-task/internal/database"
	"repartners-home-task/internal/models"
)

// PackSizeService provides business logic for pack size management
type PackSizeService struct {
	db *database.Database
}

// NewPackSizeService creates a new pack size service with database dependency
func NewPackSizeService(db *database.Database) *PackSizeService {
	return &PackSizeService{db: db}
}

// GetAllPackSizes retrieves all pack sizes from the database
func (s *PackSizeService) GetAllPackSizes() ([]models.PackSize, error) {
	return s.db.GetAllPackSizes()
}

// ReplacePackSizes replaces all pack sizes in the database with new ones
func (s *PackSizeService) ReplacePackSizes(sizes []int) error {
	return s.db.ReplacePackSizes(sizes)
}
