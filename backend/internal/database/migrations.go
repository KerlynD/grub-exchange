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

	// Create MARKET system user for market maker trades (ignore if already exists)
	_, _ = db.Exec(
		`INSERT INTO users (username, email, password_hash, ticker, bio, current_share_price, shares_outstanding)
		 VALUES ('MARKET', 'market@system', '', 'MARKET', 'Automated market maker', 0, 0)
		 ON CONFLICT (username) DO NOTHING`)

	log.Println("Database migrations completed")
	return nil
}
