package database

import (
	"database/sql"
	"fmt"
	"log"

	"repartners-home-task/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// database struct wraps the sql database connection
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection and initializes tables
func NewDatabase(dbPath string) (*Database, error) {
	// open sqlite database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// verify database connection is alive
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// create database instance
	database := &Database{db: db}

	// initialize database schema
	if err := database.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// createTables creates the pack_sizes table if it doesn't exist
func (d *Database) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS pack_sizes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		size INTEGER NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := d.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create pack_sizes table: %w", err)
	}

	return nil
}

// GetAllPackSizes retrieves all pack sizes from the database, sorted by size
func (d *Database) GetAllPackSizes() ([]models.PackSize, error) {
	// query all pack sizes ordered by size
	query := "SELECT id, size, created_at, updated_at FROM pack_sizes ORDER BY size ASC"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query pack sizes: %w", err)
	}
	defer rows.Close()

	// iterate through results and build slice
	var packSizes []models.PackSize
	for rows.Next() {
		var packSize models.PackSize
		if err := rows.Scan(&packSize.ID, &packSize.Size, &packSize.CreatedAt, &packSize.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan pack size: %w", err)
		}
		packSizes = append(packSizes, packSize)
	}

	// check for any iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating pack sizes: %w", err)
	}

	return packSizes, nil
}

// ReplacePackSizes replaces all pack sizes in the database with new ones
func (d *Database) ReplacePackSizes(sizes []int) error {
	// begin transaction for atomic operation
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// delete all existing pack sizes
	_, err = tx.Exec("DELETE FROM pack_sizes")
	if err != nil {
		return fmt.Errorf("failed to delete existing pack sizes: %w", err)
	}

	// prepare insert statement for efficiency
	stmt, err := tx.Prepare("INSERT INTO pack_sizes (size) VALUES (?)")
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	// insert each new pack size
	for _, size := range sizes {
		_, err = stmt.Exec(size)
		if err != nil {
			return fmt.Errorf("failed to insert pack size %d: %w", size, err)
		}
	}

	// commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// log successful operation
	log.Printf("successfully replaced pack sizes with %d new sizes", len(sizes))
	return nil
}
