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

	// Seed achievement definitions
	achievements := []struct {
		ID, Name, Description, Icon string
	}{
		{"first_trade", "First Trade", "Execute your first buy or sell trade", "🎯"},
		{"diamond_hands", "Diamond Hands", "Hold a stock for 30+ days", "💎"},
		{"day_trader", "Day Trader", "Execute 10 trades in a single day", "⚡"},
		{"whale", "Whale", "Portfolio worth 10,000+ Grub", "🐋"},
	}
	for _, a := range achievements {
		_, _ = db.Exec(
			`INSERT INTO achievements (id, name, description, icon) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`,
			a.ID, a.Name, a.Description, a.Icon,
		)
	}

	log.Println("Database migrations completed")
	return nil
}
