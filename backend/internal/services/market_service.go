package services

import (
	"fmt"
	"grub-exchange/internal/models"
	"grub-exchange/internal/repository"
	"log"
	"sort"
	"time"
)

type MarketService struct {
	userRepo      *repository.UserRepo
	balanceRepo   *repository.BalanceRepo
	portfolioRepo *repository.PortfolioRepo
	txnRepo       *repository.TransactionRepo
	snapshotRepo  *repository.MarketSnapshotRepo
	notifRepo     *repository.NotificationRepo
}

func NewMarketService(
	userRepo *repository.UserRepo,
	balanceRepo *repository.BalanceRepo,
	portfolioRepo *repository.PortfolioRepo,
	txnRepo *repository.TransactionRepo,
	snapshotRepo *repository.MarketSnapshotRepo,
	notifRepo *repository.NotificationRepo,
) *MarketService {
	return &MarketService{
		userRepo:      userRepo,
		balanceRepo:   balanceRepo,
		portfolioRepo: portfolioRepo,
		txnRepo:       txnRepo,
		snapshotRepo:  snapshotRepo,
		notifRepo:     notifRepo,
	}
}

func (s *MarketService) GetAllStocks() ([]models.StockListItem, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var stocks []models.StockListItem
	for _, u := range users {
		price24hAgo, _ := s.txnRepo.GetPriceAt(u.ID, time.Now().Add(-24*time.Hour))
		changePercent := 0.0
		if price24hAgo > 0 {
			changePercent = ((u.CurrentSharePrice - price24hAgo) / price24hAgo) * 100
		}

		// Get sparkline data (last 20 price points)
		history, _ := s.txnRepo.GetPriceHistory(u.ID, time.Now().Add(-7*24*time.Hour))
		var sparkline []float64
		for _, h := range history {
			sparkline = append(sparkline, h.Price)
		}
		if len(sparkline) == 0 {
			sparkline = []float64{u.CurrentSharePrice}
		}

		stocks = append(stocks, models.StockListItem{
			ID:                u.ID,
			Username:          u.Username,
			Ticker:            u.Ticker,
			CurrentSharePrice: u.CurrentSharePrice,
			Change24hPercent:  changePercent,
			SparklineData:     sparkline,
		})
	}

	return stocks, nil
}

func (s *MarketService) GetMarketOverview() (*models.MarketOverview, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var totalMarketCap float64
	for _, u := range users {
		totalMarketCap += float64(u.SharesOutstanding) * u.CurrentSharePrice
	}

	// Get all balances to compute total grub in circulation vs invested
	balances, err := s.balanceRepo.GetAllBalances()
	if err != nil {
		return nil, err
	}

	var totalCash float64
	for _, b := range balances {
		totalCash += b.GrubBalance
	}

	// Total invested = total market cap of all holdings (approximate)
	holdings, err := s.portfolioRepo.GetAllHoldings()
	if err != nil {
		return nil, err
	}
	var totalInvested float64
	for _, h := range holdings {
		stockUser, err := s.userRepo.GetByID(h.StockUserID)
		if err != nil {
			continue
		}
		totalInvested += h.NumShares * stockUser.CurrentSharePrice
	}

	totalGrub := totalCash + totalInvested
	investedPercent := 0.0
	if totalGrub > 0 {
		investedPercent = (totalInvested / totalGrub) * 100
	}

	// Get historical snapshots (last 30 days)
	history, _ := s.snapshotRepo.GetSince(time.Now().Add(-30 * 24 * time.Hour))
	if history == nil {
		history = []models.MarketSnapshot{}
	}

	return &models.MarketOverview{
		TotalMarketCap:  totalMarketCap,
		TotalGrub:       totalGrub,
		TotalInvested:   totalInvested,
		TotalCash:       totalCash,
		InvestedPercent: investedPercent,
		TotalStocks:     len(users),
		History:         history,
	}, nil
}

// RecordMarketSnapshot computes current market stats and saves a snapshot
func (s *MarketService) RecordMarketSnapshot() {
	overview, err := s.GetMarketOverview()
	if err != nil {
		log.Printf("Error computing market snapshot: %v", err)
		return
	}

	snap := &models.MarketSnapshot{
		TotalMarketCap: overview.TotalMarketCap,
		TotalInvested:  overview.TotalInvested,
		TotalCash:      overview.TotalCash,
		TotalGrub:      overview.TotalGrub,
	}

	if err := s.snapshotRepo.Record(snap); err != nil {
		log.Printf("Error recording market snapshot: %v", err)
	}
}

func (s *MarketService) GetStockDetail(ticker string) (*models.StockDetail, error) {
	user, err := s.userRepo.GetByTicker(ticker)
	if err != nil {
		return nil, err
	}

	history, err := s.txnRepo.GetPriceHistory(user.ID, time.Now().Add(-30*24*time.Hour))
	if err != nil {
		history = []models.PriceHistory{}
	}

	trades, err := s.txnRepo.GetByStock(user.ID, 10)
	if err != nil {
		trades = []models.TransactionWithDetails{}
	}

	price24hAgo, _ := s.txnRepo.GetPriceAt(user.ID, time.Now().Add(-24*time.Hour))
	change24h := user.CurrentSharePrice - price24hAgo
	changePercent := 0.0
	if price24hAgo > 0 {
		changePercent = (change24h / price24hAgo) * 100
	}

	volume, _ := s.txnRepo.GetVolume24h(user.ID)
	ath, atl, _ := s.txnRepo.GetAllTimePriceRange(user.ID)

	return &models.StockDetail{
		User:             *user,
		PriceHistory:     history,
		RecentTrades:     trades,
		Change24h:        change24h,
		Change24hPercent: changePercent,
		MarketCap:        float64(user.SharesOutstanding) * user.CurrentSharePrice,
		Volume24h:        volume,
		AllTimeHigh:      ath,
		AllTimeLow:       atl,
	}, nil
}

func (s *MarketService) GetLeaderboard() (*models.LeaderboardData, error) {
	// --- Batch-load all data upfront to avoid N+1 queries ---
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Build user lookup map by ID
	userMap := make(map[int]models.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	// Batch: all prices 24h ago (single query)
	prices24hAgo, _ := s.txnRepo.GetPricesAtBatch(time.Now().Add(-24 * time.Hour))

	// Batch: all balances (single query)
	allBalances, _ := s.balanceRepo.GetAllBalances()
	balanceMap := make(map[int]float64, len(allBalances))
	for _, b := range allBalances {
		balanceMap[b.UserID] = b.GrubBalance
	}

	// Batch: all holdings (single query)
	allHoldings, _ := s.portfolioRepo.GetAllHoldings()
	// Group holdings by owner
	holdingsByOwner := make(map[int][]models.Portfolio)
	for _, h := range allHoldings {
		holdingsByOwner[h.OwnerID] = append(holdingsByOwner[h.OwnerID], h)
	}

	// --- 1. Most Valuable Stocks ---
	sortedByPrice := make([]models.User, len(users))
	copy(sortedByPrice, users)
	sort.Slice(sortedByPrice, func(i, j int) bool {
		return sortedByPrice[i].CurrentSharePrice > sortedByPrice[j].CurrentSharePrice
	})

	var mostValuable []models.LeaderboardEntry
	for i, u := range sortedByPrice {
		if i >= 10 {
			break
		}
		mostValuable = append(mostValuable, models.LeaderboardEntry{
			Rank:     i + 1,
			UserID:   u.ID,
			Username: u.Username,
			Ticker:   u.Ticker,
			Value:    u.CurrentSharePrice,
		})
	}

	// --- 2. Biggest Gainers and Losers (24h) ---
	type userChange struct {
		user   models.User
		change float64
	}
	var changes []userChange
	for _, u := range users {
		price24hAgo := 10.0
		if p, ok := prices24hAgo[u.ID]; ok {
			price24hAgo = p
		}
		changePercent := 0.0
		if price24hAgo > 0 {
			changePercent = ((u.CurrentSharePrice - price24hAgo) / price24hAgo) * 100
		}
		changes = append(changes, userChange{user: u, change: changePercent})
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].change > changes[j].change
	})

	var gainers []models.LeaderboardEntry
	for i, c := range changes {
		if i >= 10 {
			break
		}
		gainers = append(gainers, models.LeaderboardEntry{
			Rank:     i + 1,
			UserID:   c.user.ID,
			Username: c.user.Username,
			Ticker:   c.user.Ticker,
			Value:    c.user.CurrentSharePrice,
			Change:   c.change,
		})
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].change < changes[j].change
	})

	var losers []models.LeaderboardEntry
	for i, c := range changes {
		if i >= 10 {
			break
		}
		losers = append(losers, models.LeaderboardEntry{
			Rank:     i + 1,
			UserID:   c.user.ID,
			Username: c.user.Username,
			Ticker:   c.user.Ticker,
			Value:    c.user.CurrentSharePrice,
			Change:   c.change,
		})
	}

	// --- 3. Richest Traders (cash + holdings using pre-loaded data) ---
	type userWealth struct {
		user       models.User
		totalValue float64
	}
	var wealthEntries []userWealth
	for _, u := range users {
		cash := balanceMap[u.ID]
		holdingsValue := 0.0
		for _, h := range holdingsByOwner[u.ID] {
			if su, ok := userMap[h.StockUserID]; ok {
				holdingsValue += h.NumShares * su.CurrentSharePrice
			}
		}
		wealthEntries = append(wealthEntries, userWealth{user: u, totalValue: cash + holdingsValue})
	}
	sort.Slice(wealthEntries, func(i, j int) bool {
		return wealthEntries[i].totalValue > wealthEntries[j].totalValue
	})

	var richest []models.LeaderboardEntry
	for i, w := range wealthEntries {
		if i >= 10 {
			break
		}
		richest = append(richest, models.LeaderboardEntry{
			Rank:     i + 1,
			UserID:   w.user.ID,
			Username: w.user.Username,
			Ticker:   w.user.Ticker,
			Value:    w.totalValue,
		})
	}

	// --- 4. Best Portfolio Performance (using pre-loaded data) ---
	var perfEntries []models.LeaderboardEntry
	for _, u := range users {
		holdings := holdingsByOwner[u.ID]
		if len(holdings) == 0 {
			continue
		}
		var totalValue, totalCost float64
		for _, h := range holdings {
			if su, ok := userMap[h.StockUserID]; ok {
				totalValue += h.NumShares * su.CurrentSharePrice
				totalCost += h.NumShares * h.AvgPurchasePrice
			}
		}
		plPercent := 0.0
		if totalCost > 0 {
			plPercent = ((totalValue - totalCost) / totalCost) * 100
		}
		perfEntries = append(perfEntries, models.LeaderboardEntry{
			UserID:   u.ID,
			Username: u.Username,
			Ticker:   u.Ticker,
			Value:    plPercent,
		})
	}
	sort.Slice(perfEntries, func(i, j int) bool {
		return perfEntries[i].Value > perfEntries[j].Value
	})
	for i := range perfEntries {
		if i >= 10 {
			perfEntries = perfEntries[:10]
			break
		}
		perfEntries[i].Rank = i + 1
	}

	return &models.LeaderboardData{
		MostValuable:    mostValuable,
		BiggestGainers:  gainers,
		BiggestLosers:   losers,
		RichestTraders:  richest,
		BestPerformance: perfEntries,
	}, nil
}

func (s *MarketService) GetRecentTransactions(limit int) ([]models.TransactionWithDetails, error) {
	return s.txnRepo.GetRecent(limit)
}

// RunDailyDecay applies price decay to stocks not traded in 24h
func (s *MarketService) RunDailyDecay() {
	since := time.Now().Add(-24 * time.Hour)
	users, err := s.userRepo.GetStocksNotTradedSince(since)
	if err != nil {
		log.Printf("Error getting stocks for decay: %v", err)
		return
	}

	for _, u := range users {
		newPrice := ApplyDecay(u.CurrentSharePrice)
		if newPrice != u.CurrentSharePrice {
			if err := s.userRepo.UpdateSharePriceNoTx(u.ID, newPrice); err != nil {
				log.Printf("Error applying decay for user %d: %v", u.ID, err)
				continue
			}
			if err := s.txnRepo.RecordPriceHistoryNoTx(u.ID, newPrice); err != nil {
				log.Printf("Error recording decay price for user %d: %v", u.ID, err)
			}
		}
	}

	log.Printf("Daily decay applied to %d stocks", len(users))
}

// RunDailyDividends pays 1% daily dividend on holdings value
func (s *MarketService) RunDailyDividends() {
	users, err := s.userRepo.GetAll()
	if err != nil {
		log.Printf("Error getting users for dividends: %v", err)
		return
	}

	for _, u := range users {
		holdings, err := s.portfolioRepo.GetByOwner(u.ID)
		if err != nil || len(holdings) == 0 {
			continue
		}

		var totalValue float64
		for _, h := range holdings {
			stockUser, err := s.userRepo.GetByID(h.StockUserID)
			if err != nil {
				continue
			}
			totalValue += h.NumShares * stockUser.CurrentSharePrice
		}

		dividend := totalValue * 0.01
		if dividend > 0.01 { // Only pay if dividend is meaningful (> 1 cent)
			if err := s.balanceRepo.UpdateBalanceNoTx(u.ID, dividend); err != nil {
				log.Printf("Error paying dividend to user %d: %v", u.ID, err)
				continue
			}
			// Send notification about dividend payment
			if s.notifRepo != nil {
				msg := fmt.Sprintf("You received %.2f Grub in dividends from your holdings!", dividend)
				_ = s.notifRepo.Create(u.ID, "dividend", msg, "", "", 0)
			}
		}
	}

	log.Println("Daily dividends paid")
}
