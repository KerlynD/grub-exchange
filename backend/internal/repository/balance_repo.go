package repository

import (
	"database/sql"
	"grub-exchange/internal/models"
	"time"
)

type BalanceRepo struct {
	db *sql.DB
}

func NewBalanceRepo(db *sql.DB) *BalanceRepo {
	return &BalanceRepo{db: db}
}

func (r *BalanceRepo) Create(userID int, initialBalance float64) error {
	_, err := r.db.Exec(
		`INSERT INTO balances (user_id, grub_balance) VALUES ($1, $2)`,
		userID, initialBalance,
	)
	return err
}

func (r *BalanceRepo) GetByUserID(userID int) (*models.Balance, error) {
	balance := &models.Balance{}
	err := r.db.QueryRow(
		`SELECT user_id, grub_balance, last_daily_claim FROM balances WHERE user_id = $1`, userID,
	).Scan(&balance.UserID, &balance.GrubBalance, &balance.LastDailyClaim)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (r *BalanceRepo) UpdateBalance(tx *sql.Tx, userID int, amount float64) error {
	_, err := tx.Exec(
		`UPDATE balances SET grub_balance = grub_balance + $1 WHERE user_id = $2`,
		amount, userID,
	)
	return err
}

func (r *BalanceRepo) UpdateBalanceNoTx(userID int, amount float64) error {
	_, err := r.db.Exec(
		`UPDATE balances SET grub_balance = grub_balance + $1 WHERE user_id = $2`,
		amount, userID,
	)
	return err
}

func (r *BalanceRepo) ClaimDailyBonus(userID int) error {
	_, err := r.db.Exec(
		`UPDATE balances SET grub_balance = grub_balance + 10, last_daily_claim = $1 WHERE user_id = $2`,
		time.Now(), userID,
	)
	return err
}

func (r *BalanceRepo) GetAllBalances() ([]models.Balance, error) {
	rows, err := r.db.Query(`SELECT user_id, grub_balance, last_daily_claim FROM balances`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []models.Balance
	for rows.Next() {
		var b models.Balance
		if err := rows.Scan(&b.UserID, &b.GrubBalance, &b.LastDailyClaim); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}

func (r *BalanceRepo) GetTopByBalance(limit int) ([]models.Balance, error) {
	rows, err := r.db.Query(
		`SELECT user_id, grub_balance, last_daily_claim FROM balances ORDER BY grub_balance DESC LIMIT $1`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []models.Balance
	for rows.Next() {
		var b models.Balance
		if err := rows.Scan(&b.UserID, &b.GrubBalance, &b.LastDailyClaim); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}
