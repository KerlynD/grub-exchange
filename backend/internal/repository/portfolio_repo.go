package repository

import (
	"database/sql"
	"grub-exchange/internal/models"
)

type PortfolioRepo struct {
	db *sql.DB
}

func NewPortfolioRepo(db *sql.DB) *PortfolioRepo {
	return &PortfolioRepo{db: db}
}

func (r *PortfolioRepo) GetByOwner(ownerID int) ([]models.Portfolio, error) {
	rows, err := r.db.Query(
		`SELECT id, owner_id, stock_user_id, num_shares, avg_purchase_price
		 FROM portfolios WHERE owner_id = $1`, ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var portfolios []models.Portfolio
	for rows.Next() {
		var p models.Portfolio
		if err := rows.Scan(&p.ID, &p.OwnerID, &p.StockUserID, &p.NumShares, &p.AvgPurchasePrice); err != nil {
			return nil, err
		}
		portfolios = append(portfolios, p)
	}
	return portfolios, nil
}

func (r *PortfolioRepo) GetHolding(ownerID, stockUserID int) (*models.Portfolio, error) {
	p := &models.Portfolio{}
	err := r.db.QueryRow(
		`SELECT id, owner_id, stock_user_id, num_shares, avg_purchase_price
		 FROM portfolios WHERE owner_id = $1 AND stock_user_id = $2`,
		ownerID, stockUserID,
	).Scan(&p.ID, &p.OwnerID, &p.StockUserID, &p.NumShares, &p.AvgPurchasePrice)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PortfolioRepo) UpsertHolding(tx *sql.Tx, ownerID, stockUserID int, numShares, avgPrice float64) error {
	_, err := tx.Exec(
		`INSERT INTO portfolios (owner_id, stock_user_id, num_shares, avg_purchase_price)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT(owner_id, stock_user_id) DO UPDATE SET
		 num_shares = $3, avg_purchase_price = $4`,
		ownerID, stockUserID, numShares, avgPrice,
	)
	return err
}

func (r *PortfolioRepo) ReduceShares(tx *sql.Tx, ownerID, stockUserID int, numShares float64) error {
	_, err := tx.Exec(
		`UPDATE portfolios SET num_shares = num_shares - $1 WHERE owner_id = $2 AND stock_user_id = $3`,
		numShares, ownerID, stockUserID,
	)
	return err
}

func (r *PortfolioRepo) DeleteHolding(tx *sql.Tx, ownerID, stockUserID int) error {
	_, err := tx.Exec(
		`DELETE FROM portfolios WHERE owner_id = $1 AND stock_user_id = $2`,
		ownerID, stockUserID,
	)
	return err
}

func (r *PortfolioRepo) GetAllHoldings() ([]models.Portfolio, error) {
	rows, err := r.db.Query(
		`SELECT id, owner_id, stock_user_id, num_shares, avg_purchase_price FROM portfolios`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var portfolios []models.Portfolio
	for rows.Next() {
		var p models.Portfolio
		if err := rows.Scan(&p.ID, &p.OwnerID, &p.StockUserID, &p.NumShares, &p.AvgPurchasePrice); err != nil {
			return nil, err
		}
		portfolios = append(portfolios, p)
	}
	return portfolios, nil
}
