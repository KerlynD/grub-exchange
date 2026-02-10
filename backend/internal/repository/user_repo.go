package repository

import (
	"database/sql"
	"grub-exchange/internal/models"
	"time"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *models.User) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO users (username, email, password_hash, ticker, bio, current_share_price, shares_outstanding, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		user.Username, user.Email, user.PasswordHash, user.Ticker, user.Bio,
		user.CurrentSharePrice, user.SharesOutstanding, time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepo) scanUser(row interface{ Scan(...interface{}) error }) (*models.User, error) {
	user := &models.User{}
	var bio sql.NullString
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Ticker, &bio, &user.CurrentSharePrice, &user.SharesOutstanding,
		&user.LastLogin, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	if bio.Valid {
		user.Bio = bio.String
	}
	return user, nil
}

const userSelectCols = `id, username, email, password_hash, ticker, bio, current_share_price, shares_outstanding, last_login, created_at`

func (r *UserRepo) GetByID(id int) (*models.User, error) {
	return r.scanUser(r.db.QueryRow(
		`SELECT `+userSelectCols+` FROM users WHERE id = $1`, id,
	))
}

func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	return r.scanUser(r.db.QueryRow(
		`SELECT `+userSelectCols+` FROM users WHERE email = $1`, email,
	))
}

func (r *UserRepo) GetByTicker(ticker string) (*models.User, error) {
	return r.scanUser(r.db.QueryRow(
		`SELECT `+userSelectCols+` FROM users WHERE ticker = $1`, ticker,
	))
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	return r.scanUser(r.db.QueryRow(
		`SELECT `+userSelectCols+` FROM users WHERE username = $1`, username,
	))
}

func (r *UserRepo) GetAll() ([]models.User, error) {
	rows, err := r.db.Query(
		`SELECT `+userSelectCols+` FROM users WHERE username != 'MARKET' ORDER BY current_share_price DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := &models.User{}
		var bio sql.NullString
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.Ticker, &bio, &user.CurrentSharePrice, &user.SharesOutstanding,
			&user.LastLogin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		if bio.Valid {
			user.Bio = bio.String
		}
		users = append(users, *user)
	}
	return users, nil
}

func (r *UserRepo) UpdateSharePrice(tx *sql.Tx, userID int, newPrice float64) error {
	_, err := tx.Exec(
		`UPDATE users SET current_share_price = $1 WHERE id = $2`,
		newPrice, userID,
	)
	return err
}

func (r *UserRepo) UpdateSharePriceNoTx(userID int, newPrice float64) error {
	_, err := r.db.Exec(
		`UPDATE users SET current_share_price = $1 WHERE id = $2`,
		newPrice, userID,
	)
	return err
}

func (r *UserRepo) UpdateBio(userID int, bio string) error {
	_, err := r.db.Exec(
		`UPDATE users SET bio = $1 WHERE id = $2`,
		bio, userID,
	)
	return err
}

func (r *UserRepo) UpdateLastLogin(userID int) error {
	_, err := r.db.Exec(
		`UPDATE users SET last_login = $1 WHERE id = $2`,
		time.Now(), userID,
	)
	return err
}

func (r *UserRepo) ExistsByEmail(email string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = $1`, email).Scan(&count)
	return count > 0, err
}

func (r *UserRepo) ExistsByUsername(username string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = $1`, username).Scan(&count)
	return count > 0, err
}

func (r *UserRepo) ExistsByTicker(ticker string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE ticker = $1`, ticker).Scan(&count)
	return count > 0, err
}

func (r *UserRepo) GetStocksNotTradedSince(since time.Time) ([]models.User, error) {
	rows, err := r.db.Query(
		`SELECT `+userSelectCols+` FROM users u
		 WHERE u.id NOT IN (
			 SELECT DISTINCT stock_user_id FROM transactions WHERE timestamp > $1
		 )`, since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		user := &models.User{}
		var bio sql.NullString
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.Ticker, &bio, &user.CurrentSharePrice, &user.SharesOutstanding,
			&user.LastLogin, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		if bio.Valid {
			user.Bio = bio.String
		}
		users = append(users, *user)
	}
	return users, nil
}

// Portfolio snapshot methods
func (r *UserRepo) SavePortfolioSnapshot(userID int, totalValue, grubBalance float64) error {
	_, err := r.db.Exec(
		`INSERT INTO portfolio_snapshots (user_id, total_value, grub_balance, timestamp) VALUES ($1, $2, $3, $4)`,
		userID, totalValue, grubBalance, time.Now(),
	)
	return err
}

func (r *UserRepo) GetPortfolioSnapshots(userID int, since time.Time) ([]models.PortfolioSnapshot, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, total_value, grub_balance, timestamp FROM portfolio_snapshots
		 WHERE user_id = $1 AND timestamp > $2 ORDER BY timestamp ASC`,
		userID, since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []models.PortfolioSnapshot
	for rows.Next() {
		var s models.PortfolioSnapshot
		if err := rows.Scan(&s.ID, &s.UserID, &s.TotalValue, &s.GrubBalance, &s.Timestamp); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}
