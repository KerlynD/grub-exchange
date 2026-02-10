package database

import (
	"database/sql"
	_ "embed"
	"log"
)

//go:embed schema.sql
var schema string

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		log.Printf("Error running migrations: %v", err)
		return err
	}

	// Add bio column if it doesn't exist (SQLite ALTER TABLE)
	runSafe(db, `ALTER TABLE users ADD COLUMN bio TEXT DEFAULT ''`)

	// Create MARKET system user for market maker trades (ignore if already exists)
	runSafe(db, `INSERT INTO users (username, email, password_hash, ticker, bio, current_share_price, shares_outstanding, created_at)
		VALUES ('MARKET', 'market@system', '', 'MARKET', 'Automated market maker', 0, 0, CURRENT_TIMESTAMP)`)

	log.Println("Database migrations completed")
	return nil
}

// runSafe executes a SQL statement and ignores errors (e.g., column already exists)
func runSafe(db *sql.DB, query string) {
	_, err := db.Exec(query)
	if err != nil {
		// Ignore "duplicate column" errors
		log.Printf("Migration note (safe to ignore): %v", err)
	}
}
