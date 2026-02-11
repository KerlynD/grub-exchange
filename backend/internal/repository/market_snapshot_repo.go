package repository

import (
	"database/sql"
	"grub-exchange/internal/models"
	"log"
	"time"
)

type MarketSnapshotRepo struct {
	db *sql.DB
}

func NewMarketSnapshotRepo(db *sql.DB) *MarketSnapshotRepo {
	return &MarketSnapshotRepo{db: db}
}

func (r *MarketSnapshotRepo) Record(snap *models.MarketSnapshot) error {
	_, err := r.db.Exec(
		`INSERT INTO market_snapshots (total_market_cap, total_invested, total_cash, total_grub, timestamp)
		 VALUES ($1, $2, $3, $4, $5)`,
		snap.TotalMarketCap, snap.TotalInvested, snap.TotalCash, snap.TotalGrub, time.Now(),
	)
	return err
}

func (r *MarketSnapshotRepo) GetSince(since time.Time) ([]models.MarketSnapshot, error) {
	rows, err := r.db.Query(
		`SELECT id, total_market_cap, total_invested, total_cash, total_grub, timestamp
		 FROM market_snapshots WHERE timestamp > $1 ORDER BY timestamp ASC`,
		since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []models.MarketSnapshot
	for rows.Next() {
		var s models.MarketSnapshot
		if err := rows.Scan(&s.ID, &s.TotalMarketCap, &s.TotalInvested, &s.TotalCash, &s.TotalGrub, &s.Timestamp); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

func (r *MarketSnapshotRepo) Count() int {
	var count int
	r.db.QueryRow(`SELECT COUNT(*) FROM market_snapshots`).Scan(&count)
	return count
}

// BackfillFromHistory generates hourly market snapshots from existing price_history
// and portfolio_snapshots data. Only runs if market_snapshots is empty.
func (r *MarketSnapshotRepo) BackfillFromHistory() {
	if r.Count() > 0 {
		return
	}

	log.Println("Backfilling market snapshots from historical data...")

	// Use a single PostgreSQL query with generate_series + DISTINCT ON
	// to compute hourly market cap from price_history, and hourly invested/cash from portfolio_snapshots.
	result, err := r.db.Exec(`
		INSERT INTO market_snapshots (total_market_cap, total_invested, total_cash, total_grub, timestamp)
		SELECT
			COALESCE(mc.total_market_cap, 0),
			COALESCE(ps.total_invested, 0),
			COALESCE(ps.total_cash, 0),
			COALESCE(ps.total_invested, 0) + COALESCE(ps.total_cash, 0),
			tb.bucket_time
		FROM generate_series(
			(SELECT MIN(timestamp) FROM price_history),
			NOW(),
			interval '1 hour'
		) AS tb(bucket_time)
		LEFT JOIN LATERAL (
			SELECT SUM(lp.price * 1000) AS total_market_cap
			FROM (
				SELECT DISTINCT ON (user_id) price
				FROM price_history
				WHERE timestamp <= tb.bucket_time
				ORDER BY user_id, timestamp DESC
			) lp
		) mc ON true
		LEFT JOIN LATERAL (
			SELECT
				SUM(ls.total_value) AS total_invested,
				SUM(ls.grub_balance) AS total_cash
			FROM (
				SELECT DISTINCT ON (user_id) total_value, grub_balance
				FROM portfolio_snapshots
				WHERE timestamp <= tb.bucket_time
				ORDER BY user_id, timestamp DESC
			) ls
		) ps ON true
		WHERE mc.total_market_cap IS NOT NULL
		ORDER BY tb.bucket_time
	`)
	if err != nil {
		log.Printf("Error backfilling market snapshots: %v", err)
		return
	}

	rows, _ := result.RowsAffected()
	log.Printf("Backfilled %d market snapshots from historical data", rows)
}
