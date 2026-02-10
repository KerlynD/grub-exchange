package services

import (
	"errors"
	"grub-exchange/internal/models"
	"grub-exchange/internal/repository"
	"grub-exchange/internal/utils"
	"strings"
)

type AuthService struct {
	userRepo    *repository.UserRepo
	balanceRepo *repository.BalanceRepo
	txnRepo     *repository.TransactionRepo
}

func NewAuthService(userRepo *repository.UserRepo, balanceRepo *repository.BalanceRepo, txnRepo *repository.TransactionRepo) *AuthService {
	return &AuthService{userRepo: userRepo, balanceRepo: balanceRepo, txnRepo: txnRepo}
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

	user := &models.User{
		Username:          req.Username,
		Email:             strings.ToLower(req.Email),
		PasswordHash:      hashedPassword,
		Ticker:            ticker,
		CurrentSharePrice: 10.0,
		SharesOutstanding: 1000,
	}

	userID, err := s.userRepo.Create(user)
	if err != nil {
		return nil, "", err
	}

	user.ID = int(userID)

	if err := s.balanceRepo.Create(user.ID, 100.0); err != nil {
		return nil, "", err
	}

	// Record initial price in history
	if err := s.txnRepo.RecordPriceHistoryNoTx(user.ID, 10.0); err != nil {
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
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
		GrubBalance:       100.0,
		CreatedAt:         user.CreatedAt,
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
