package services

import (
	"grub-exchange/internal/models"
	"grub-exchange/internal/repository"
	"log"
)

type AchievementService struct {
	achievementRepo *repository.AchievementRepo
	balanceRepo     *repository.BalanceRepo
	portfolioRepo   *repository.PortfolioRepo
	userRepo        *repository.UserRepo
}

func NewAchievementService(
	achievementRepo *repository.AchievementRepo,
	balanceRepo *repository.BalanceRepo,
	portfolioRepo *repository.PortfolioRepo,
	userRepo *repository.UserRepo,
) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		balanceRepo:     balanceRepo,
		portfolioRepo:   portfolioRepo,
		userRepo:        userRepo,
	}
}

// CheckAfterTrade checks and awards trade-related achievements
func (s *AchievementService) CheckAfterTrade(userID int) []models.UserAchievement {
	var newlyEarned []models.UserAchievement

	// First Trade
	if earned := s.checkFirstTrade(userID); earned != nil {
		newlyEarned = append(newlyEarned, *earned)
	}

	// Day Trader (10 trades in one day)
	if earned := s.checkDayTrader(userID); earned != nil {
		newlyEarned = append(newlyEarned, *earned)
	}

	// Whale (portfolio worth 10,000+ Grub)
	if earned := s.checkWhale(userID); earned != nil {
		newlyEarned = append(newlyEarned, *earned)
	}

	return newlyEarned
}

// CheckDiamondHands checks if a user has held a stock for 30+ days (called periodically)
func (s *AchievementService) CheckDiamondHands(userID int) *models.UserAchievement {
	has, _ := s.achievementRepo.HasAchievement(userID, "diamond_hands")
	if has {
		return nil
	}

	days, err := s.achievementRepo.GetOldestHoldingDays(userID)
	if err != nil || days < 30 {
		return nil
	}

	if err := s.achievementRepo.Award(userID, "diamond_hands"); err != nil {
		log.Printf("Error awarding diamond_hands to user %d: %v", userID, err)
		return nil
	}

	return &models.UserAchievement{
		UserID:        userID,
		AchievementID: "diamond_hands",
		Name:          "Diamond Hands",
		Description:   "Hold a stock for 30+ days",
		Icon:          "💎",
	}
}

func (s *AchievementService) checkFirstTrade(userID int) *models.UserAchievement {
	has, _ := s.achievementRepo.HasAchievement(userID, "first_trade")
	if has {
		return nil
	}

	if err := s.achievementRepo.Award(userID, "first_trade"); err != nil {
		log.Printf("Error awarding first_trade to user %d: %v", userID, err)
		return nil
	}

	return &models.UserAchievement{
		UserID:        userID,
		AchievementID: "first_trade",
		Name:          "First Trade",
		Description:   "Execute your first buy or sell trade",
		Icon:          "🎯",
	}
}

func (s *AchievementService) checkDayTrader(userID int) *models.UserAchievement {
	has, _ := s.achievementRepo.HasAchievement(userID, "day_trader")
	if has {
		return nil
	}

	count, err := s.achievementRepo.GetTradeCountToday(userID)
	if err != nil || count < 10 {
		return nil
	}

	if err := s.achievementRepo.Award(userID, "day_trader"); err != nil {
		log.Printf("Error awarding day_trader to user %d: %v", userID, err)
		return nil
	}

	return &models.UserAchievement{
		UserID:        userID,
		AchievementID: "day_trader",
		Name:          "Day Trader",
		Description:   "Execute 10 trades in a single day",
		Icon:          "⚡",
	}
}

func (s *AchievementService) checkWhale(userID int) *models.UserAchievement {
	has, _ := s.achievementRepo.HasAchievement(userID, "whale")
	if has {
		return nil
	}

	balance, err := s.balanceRepo.GetByUserID(userID)
	if err != nil {
		return nil
	}

	holdings, err := s.portfolioRepo.GetByOwner(userID)
	if err != nil {
		return nil
	}

	totalValue := balance.GrubBalance
	for _, h := range holdings {
		stockUser, err := s.userRepo.GetByID(h.StockUserID)
		if err != nil {
			continue
		}
		totalValue += h.NumShares * stockUser.CurrentSharePrice
	}

	if totalValue < 10000 {
		return nil
	}

	if err := s.achievementRepo.Award(userID, "whale"); err != nil {
		log.Printf("Error awarding whale to user %d: %v", userID, err)
		return nil
	}

	return &models.UserAchievement{
		UserID:        userID,
		AchievementID: "whale",
		Name:          "Whale",
		Description:   "Portfolio worth 10,000+ Grub",
		Icon:          "🐋",
	}
}

func (s *AchievementService) GetUserAchievements(userID int) ([]models.UserAchievement, error) {
	return s.achievementRepo.GetByUser(userID)
}

func (s *AchievementService) GetAllAchievements() ([]models.Achievement, error) {
	return s.achievementRepo.GetAll()
}
