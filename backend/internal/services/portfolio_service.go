package services

import (
	"errors"
	"grub-exchange/internal/models"
	"grub-exchange/internal/repository"
	"time"
)

type PortfolioService struct {
	userRepo      *repository.UserRepo
	balanceRepo   *repository.BalanceRepo
	portfolioRepo *repository.PortfolioRepo
	txnRepo       *repository.TransactionRepo
}

func NewPortfolioService(
	userRepo *repository.UserRepo,
	balanceRepo *repository.BalanceRepo,
	portfolioRepo *repository.PortfolioRepo,
	txnRepo *repository.TransactionRepo,
) *PortfolioService {
	return &PortfolioService{
		userRepo:      userRepo,
		balanceRepo:   balanceRepo,
		portfolioRepo: portfolioRepo,
		txnRepo:       txnRepo,
	}
}

func (s *PortfolioService) GetUserPortfolio(userID int) (*models.PortfolioResponse, error) {
	balance, err := s.balanceRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	holdings, err := s.portfolioRepo.GetByOwner(userID)
	if err != nil {
		return nil, err
	}

	var portfolioHoldings []models.PortfolioHolding
	var totalValue float64
	var totalCost float64

	for _, h := range holdings {
		stockUser, err := s.userRepo.GetByID(h.StockUserID)
		if err != nil {
			continue
		}

		currentValue := h.NumShares * stockUser.CurrentSharePrice
		costBasis := h.NumShares * h.AvgPurchasePrice
		pl := currentValue - costBasis
		plPercent := 0.0
		if costBasis > 0 {
			plPercent = (pl / costBasis) * 100
		}

		portfolioHoldings = append(portfolioHoldings, models.PortfolioHolding{
			Ticker:            stockUser.Ticker,
			Username:          stockUser.Username,
			StockUserID:       stockUser.ID,
			NumShares:         h.NumShares,
			AvgPurchasePrice:  h.AvgPurchasePrice,
			CurrentPrice:      stockUser.CurrentSharePrice,
			TotalValue:        currentValue,
			ProfitLoss:        pl,
			ProfitLossPercent: plPercent,
		})

		totalValue += currentValue
		totalCost += costBasis
	}

	totalPL := totalValue - totalCost
	totalPLPercent := 0.0
	if totalCost > 0 {
		totalPLPercent = (totalPL / totalCost) * 100
	}

	canClaim := true
	var lastClaimStr *string
	if balance.LastDailyClaim != nil {
		canClaim = time.Since(*balance.LastDailyClaim) >= 24*time.Hour
		formatted := balance.LastDailyClaim.Format(time.RFC3339)
		lastClaimStr = &formatted
	}

	return &models.PortfolioResponse{
		GrubBalance:    balance.GrubBalance,
		TotalValue:     totalValue,
		TotalPL:        totalPL,
		TotalPLPercent: totalPLPercent,
		Holdings:       portfolioHoldings,
		CanClaimDaily:  canClaim,
		LastDailyClaim: lastClaimStr,
	}, nil
}

func (s *PortfolioService) ClaimDailyBonus(userID int) (float64, error) {
	balance, err := s.balanceRepo.GetByUserID(userID)
	if err != nil {
		return 0, err
	}

	if balance.LastDailyClaim != nil {
		if time.Since(*balance.LastDailyClaim) < 24*time.Hour {
			return 0, errors.New("daily bonus already claimed, try again later")
		}
	}

	// Get user's stock price for the bonus calculation
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return 0, err
	}

	// Daily bonus: 20 Grub base + 5% of current stock price
	baseBonus := 20.0
	stockBonus := user.CurrentSharePrice * 0.05
	totalBonus := baseBonus + stockBonus

	if err := s.balanceRepo.ClaimDailyBonus(userID, totalBonus); err != nil {
		return 0, err
	}

	return balance.GrubBalance + totalBonus, nil
}

func (s *PortfolioService) GetTransactionHistory(userID int) ([]models.TransactionWithDetails, error) {
	return s.txnRepo.GetByUser(userID, 50)
}
