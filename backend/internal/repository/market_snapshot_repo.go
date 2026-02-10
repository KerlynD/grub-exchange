package repository

import (
	"database/sql"
	"grub-exchange/internal/models"
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
