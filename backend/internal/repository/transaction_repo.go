package repository

import (
	"database/sql"
	"grub-exchange/internal/models"
	"time"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(tx *sql.Tx, t *models.Transaction) error {
	_, err := tx.Exec(
		`INSERT INTO transactions (buyer_id, stock_user_id, transaction_type, num_shares, price_per_share, total_grub, timestamp)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		t.BuyerID, t.StockUserID, t.TransactionType, t.NumShares, t.PricePerShare, t.TotalGrub, time.Now(),
	)
	return err
}

func (r *TransactionRepo) CreateNoTx(t *models.Transaction) error {
	_, err := r.db.Exec(
		`INSERT INTO transactions (buyer_id, stock_user_id, transaction_type, num_shares, price_per_share, total_grub, timestamp)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		t.BuyerID, t.StockUserID, t.TransactionType, t.NumShares, t.PricePerShare, t.TotalGrub, time.Now(),
	)
	return err
}

func (r *TransactionRepo) GetByUser(userID int, limit int) ([]models.TransactionWithDetails, error) {
	rows, err := r.db.Query(
		`SELECT t.id, u1.username, u2.ticker, t.transaction_type, t.num_shares, t.price_per_share, t.total_grub, t.timestamp
		 FROM transactions t
		 JOIN users u1 ON t.buyer_id = u1.id
		 JOIN users u2 ON t.stock_user_id = u2.id
		 WHERE t.buyer_id = $1
		 ORDER BY t.timestamp DESC LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []models.TransactionWithDetails
	for rows.Next() {
		var t models.TransactionWithDetails
		if err := rows.Scan(&t.ID, &t.BuyerUsername, &t.StockTicker, &t.TransactionType,
			&t.NumShares, &t.PricePerShare, &t.TotalGrub, &t.Timestamp); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, nil
}

func (r *TransactionRepo) GetByStock(stockUserID int, limit int) ([]models.TransactionWithDetails, error) {
	rows, err := r.db.Query(
		`SELECT t.id, u1.username, u2.ticker, t.transaction_type, t.num_shares, t.price_per_share, t.total_grub, t.timestamp
		 FROM transactions t
		 JOIN users u1 ON t.buyer_id = u1.id
		 JOIN users u2 ON t.stock_user_id = u2.id
		 WHERE t.stock_user_id = $1
		 ORDER BY t.timestamp DESC LIMIT $2`,
		stockUserID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []models.TransactionWithDetails
	for rows.Next() {
		var t models.TransactionWithDetails
		if err := rows.Scan(&t.ID, &t.BuyerUsername, &t.StockTicker, &t.TransactionType,
			&t.NumShares, &t.PricePerShare, &t.TotalGrub, &t.Timestamp); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, nil
}

func (r *TransactionRepo) GetRecent(limit int) ([]models.TransactionWithDetails, error) {
	rows, err := r.db.Query(
		`SELECT t.id, u1.username, u2.ticker, t.transaction_type, t.num_shares, t.price_per_share, t.total_grub, t.timestamp
		 FROM transactions t
		 JOIN users u1 ON t.buyer_id = u1.id
		 JOIN users u2 ON t.stock_user_id = u2.id
		 ORDER BY t.timestamp DESC LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []models.TransactionWithDetails
	for rows.Next() {
		var t models.TransactionWithDetails
		if err := rows.Scan(&t.ID, &t.BuyerUsername, &t.StockTicker, &t.TransactionType,
			&t.NumShares, &t.PricePerShare, &t.TotalGrub, &t.Timestamp); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, nil
}

func (r *TransactionRepo) GetVolume24h(stockUserID int) (float64, error) {
	var volume sql.NullFloat64
	err := r.db.QueryRow(
		`SELECT SUM(num_shares) FROM transactions
		 WHERE stock_user_id = $1 AND timestamp > $2`,
		stockUserID, time.Now().Add(-24*time.Hour),
	).Scan(&volume)
	if err != nil {
		return 0, err
	}
	if volume.Valid {
		return volume.Float64, nil
	}
	return 0, nil
}

func (r *TransactionRepo) RecordPriceHistory(tx *sql.Tx, userID int, price float64) error {
	_, err := tx.Exec(
		`INSERT INTO price_history (user_id, price, timestamp) VALUES ($1, $2, $3)`,
		userID, price, time.Now(),
	)
	return err
}

func (r *TransactionRepo) RecordPriceHistoryNoTx(userID int, price float64) error {
	_, err := r.db.Exec(
		`INSERT INTO price_history (user_id, price, timestamp) VALUES ($1, $2, $3)`,
		userID, price, time.Now(),
	)
	return err
}

func (r *TransactionRepo) GetPriceHistory(userID int, since time.Time) ([]models.PriceHistory, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, price, timestamp FROM price_history
		 WHERE user_id = $1 AND timestamp > $2
		 ORDER BY timestamp ASC`,
		userID, since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.PriceHistory
	for rows.Next() {
		var ph models.PriceHistory
		if err := rows.Scan(&ph.ID, &ph.UserID, &ph.Price, &ph.Timestamp); err != nil {
			return nil, err
		}
		history = append(history, ph)
	}
	return history, nil
}

func (r *TransactionRepo) GetAllTimePriceRange(userID int) (high float64, low float64, err error) {
	err = r.db.QueryRow(
		`SELECT COALESCE(MAX(price), 10.0), COALESCE(MIN(price), 10.0) FROM price_history WHERE user_id = $1`,
		userID,
	).Scan(&high, &low)
	return
}

func (r *TransactionRepo) GetPriceAt(userID int, at time.Time) (float64, error) {
	var price float64
	err := r.db.QueryRow(
		`SELECT price FROM price_history WHERE user_id = $1 AND timestamp <= $2 ORDER BY timestamp DESC LIMIT 1`,
		userID, at,
	).Scan(&price)
	if err == sql.ErrNoRows {
		return 10.0, nil
	}
	return price, err
}

// GetPricesAtBatch returns the price for each user_id at a given time in a single query.
// Uses DISTINCT ON to get the latest price before or at the timestamp per user.
func (r *TransactionRepo) GetPricesAtBatch(at time.Time) (map[int]float64, error) {
	rows, err := r.db.Query(
		`SELECT DISTINCT ON (user_id) user_id, price
		 FROM price_history
		 WHERE timestamp <= $1
		 ORDER BY user_id, timestamp DESC`,
		at,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prices := make(map[int]float64)
	for rows.Next() {
		var uid int
		var price float64
		if err := rows.Scan(&uid, &price); err != nil {
			return nil, err
		}
		prices[uid] = price
	}
	return prices, nil
}
