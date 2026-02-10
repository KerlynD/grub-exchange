package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost:5432/grub_exchange?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Hash a password for the fake user
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	type pricePoint struct {
		time  time.Time
		price float64
	}

	// Get an existing user to use as "buyer" for fake transactions
	var buyerID int
	err = db.QueryRow(`SELECT id FROM users WHERE username != 'MARKET' LIMIT 1`).Scan(&buyerID)
	if err != nil {
		// Create a placeholder buyer
		err = db.QueryRow(
			`INSERT INTO users (username, email, password_hash, ticker, current_share_price, shares_outstanding, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)
			 ON CONFLICT (username) DO UPDATE SET username = EXCLUDED.username
			 RETURNING id`,
			"trader_bot", "bot@grubexchange.com", string(hash), "BOT", 10.0, 1000, time.Now(),
		).Scan(&buyerID)
		if err != nil {
			log.Fatal("Error creating placeholder buyer: ", err)
		}
		_, _ = db.Exec(`INSERT INTO balances (user_id, grub_balance) VALUES ($1, $2) ON CONFLICT DO NOTHING`, buyerID, 10000.0)
	}

	// =========================================
	// Seed Ivan — slow bleed from 10 to 0.5 since Jan 1
	// =========================================

	var ivanID int
	err = db.QueryRow(
		`INSERT INTO users (username, email, password_hash, ticker, current_share_price, shares_outstanding, last_login, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 ON CONFLICT (username) DO UPDATE SET current_share_price = EXCLUDED.current_share_price
		 RETURNING id`,
		"ivan", "ivan@grubexchange.com", string(hash), "IVAN",
		0.50, 1000, time.Now(), time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	).Scan(&ivanID)
	if err != nil {
		log.Fatal("Error inserting Ivan: ", err)
	}

	// Clean old data if re-seeding
	db.Exec(`DELETE FROM price_history WHERE user_id = $1`, ivanID)
	db.Exec(`DELETE FROM transactions WHERE stock_user_id = $1`, ivanID)

	// Insert balance for Ivan
	_, _ = db.Exec(`INSERT INTO balances (user_id, grub_balance) VALUES ($1, $2) ON CONFLICT DO NOTHING`, ivanID, 100.0)

	// Generate price history: steady decline 10 -> 0.5 from Jan 1 to Feb 9
	ivanStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	ivanEnd := time.Date(2026, 2, 9, 23, 0, 0, 0, time.UTC)

	var ivanPoints []pricePoint
	ivanDays := ivanEnd.Sub(ivanStart).Hours() / 24
	totalIvanPts := int(ivanDays * 3)

	for i := 0; i <= totalIvanPts; i++ {
		t := ivanStart.Add(time.Duration(float64(i)/float64(totalIvanPts)*ivanEnd.Sub(ivanStart).Hours()) * time.Hour)
		progress := float64(i) / float64(totalIvanPts)

		// Exponential decay from 10.0 to 0.5
		basePrice := 10.0 * math.Exp(-3.0*progress)

		// Small noise
		noise := math.Sin(float64(i)*0.9) * basePrice * 0.04
		smallWiggle := math.Cos(float64(i)*2.3) * basePrice * 0.02

		price := basePrice + noise + smallWiggle
		price = math.Max(0.10, math.Round(price*100)/100)

		if progress > 0.97 {
			price = 0.50 + (1.0-progress)*3.0 + smallWiggle*0.1
			price = math.Max(0.10, math.Round(price*100)/100)
		}

		ivanPoints = append(ivanPoints, pricePoint{t, price})
	}

	ivanPoints = append(ivanPoints, pricePoint{ivanEnd, 0.50})

	// Insert Ivan's price history
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(`INSERT INTO price_history (user_id, price, timestamp) VALUES ($1, $2, $3)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, p := range ivanPoints {
		_, err := stmt.Exec(ivanID, p.price, p.time)
		if err != nil {
			log.Printf("Error inserting Ivan price point: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Error committing Ivan prices: ", err)
	}

	// Sparse transactions
	ivanTxns := []struct {
		txnType string
		shares  float64
		price   float64
		time    time.Time
	}{
		{"BUY", 3, 9.80, time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)},
		{"SELL", 3, 8.50, time.Date(2026, 1, 5, 14, 0, 0, 0, time.UTC)},
		{"BUY", 2, 7.20, time.Date(2026, 1, 9, 11, 30, 0, 0, time.UTC)},
		{"SELL", 2, 5.10, time.Date(2026, 1, 15, 16, 0, 0, 0, time.UTC)},
		{"SELL", 1, 2.80, time.Date(2026, 1, 25, 9, 0, 0, 0, time.UTC)},
		{"SELL", 1, 1.20, time.Date(2026, 2, 3, 12, 0, 0, 0, time.UTC)},
	}

	for _, t := range ivanTxns {
		_, err := db.Exec(
			`INSERT INTO transactions (buyer_id, stock_user_id, transaction_type, num_shares, price_per_share, total_grub, timestamp)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			buyerID, ivanID, t.txnType, t.shares, t.price, t.shares*t.price, t.time,
		)
		if err != nil {
			log.Printf("Error inserting Ivan transaction: %v", err)
		}
	}

	fmt.Printf("Seeded user 'ivan' (ID: %d) with ticker IVAN\n", ivanID)
	fmt.Printf("Current price: 0.50 Grub (down from 10.00 since Jan 1)\n")
	fmt.Printf("Inserted %d price history points\n", len(ivanPoints))
	fmt.Printf("Inserted %d fake transactions\n", len(ivanTxns))
	fmt.Println("\nDone!")
}
