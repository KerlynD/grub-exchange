package services

import (
	"database/sql"
	"errors"
	"grub-exchange/internal/models"
	"grub-exchange/internal/repository"
	"grub-exchange/internal/utils"
	"strings"
	"time"
)

type AuthService struct {
	db          *sql.DB
	userRepo    *repository.UserRepo
	balanceRepo *repository.BalanceRepo
	txnRepo     *repository.TransactionRepo
}

func NewAuthService(db *sql.DB, userRepo *repository.UserRepo, balanceRepo *repository.BalanceRepo, txnRepo *repository.TransactionRepo) *AuthService {
	return &AuthService{db: db, userRepo: userRepo, balanceRepo: balanceRepo, txnRepo: txnRepo}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.UserResponse, string, error) {
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", errors.New("email already registered")
	}

	exists, err = s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", errors.New("username already taken")
	}

	ticker := utils.SanitizeTicker(req.FirstName)
	if !utils.ValidateTicker(ticker) {
		return nil, "", errors.New("invalid first name for ticker")
	}

	exists, err = s.userRepo.ExistsByTicker(ticker)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", errors.New("ticker already exists, try a different name")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	// Atomic registration: create user, balance, and initial price history in one transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()

	var userID int64
	err = tx.QueryRow(
		`INSERT INTO users (username, email, password_hash, ticker, bio, current_share_price, shares_outstanding, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		req.Username, strings.ToLower(req.Email), hashedPassword, ticker, "", 10.0, 1000, time.Now(),
	).Scan(&userID)
	if err != nil {
		return nil, "", err
	}

	_, err = tx.Exec(
		`INSERT INTO balances (user_id, grub_balance) VALUES ($1, $2)`,
		userID, 100.0,
	)
	if err != nil {
		return nil, "", err
	}

	_, err = tx.Exec(
		`INSERT INTO price_history (user_id, price, timestamp) VALUES ($1, $2, $3)`,
		userID, 10.0, time.Now(),
	)
	if err != nil {
		return nil, "", err
	}

	if err := tx.Commit(); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(int(userID), req.Username)
	if err != nil {
		return nil, "", err
	}

	resp := &models.UserResponse{
		ID:                int(userID),
		Username:          req.Username,
		Email:             strings.ToLower(req.Email),
		Ticker:            ticker,
		Bio:               "",
		CurrentSharePrice: 10.0,
		SharesOutstanding: 1000,
		GrubBalance:       100.0,
	}

	return resp, token, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.UserResponse, string, error) {
	user, err := s.userRepo.GetByEmail(strings.ToLower(req.Email))
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return nil, "", errors.New("invalid email or password")
	}

	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	balance, err := s.balanceRepo.GetByUserID(user.ID)
	if err != nil {
		return nil, "", err
	}

	resp := &models.UserResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		Ticker:            user.Ticker,
		Bio:               user.Bio,
		CurrentSharePrice: user.CurrentSharePrice,
		SharesOutstanding: user.SharesOutstanding,
		GrubBalance:       balance.GrubBalance,
		LastLogin:         user.LastLogin,
		CreatedAt:         user.CreatedAt,
	}

	return resp, token, nil
}

func (s *AuthService) GetMe(userID int) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	balance, err := s.balanceRepo.GetByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	resp := &models.UserResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		Ticker:            user.Ticker,
		Bio:               user.Bio,
		CurrentSharePrice: user.CurrentSharePrice,
		SharesOutstanding: user.SharesOutstanding,
		GrubBalance:       balance.GrubBalance,
		LastLogin:         user.LastLogin,
		CreatedAt:         user.CreatedAt,
	}

	return resp, nil
}
