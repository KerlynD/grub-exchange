package services

import (
	"grub-exchange/internal/models"
	"grub-exchange/internal/repository"
	"log"
	"math"
	"math/rand"
	"time"
)

type MarketMaker struct {
	userRepo      *repository.UserRepo
	balanceRepo   *repository.BalanceRepo
	portfolioRepo *repository.PortfolioRepo
	txnRepo       *repository.TransactionRepo
	marketUserID  int
}

func NewMarketMaker(
	userRepo *repository.UserRepo,
	balanceRepo *repository.BalanceRepo,
	portfolioRepo *repository.PortfolioRepo,
	txnRepo *repository.TransactionRepo,
) *MarketMaker {
	return &MarketMaker{
		userRepo:      userRepo,
		balanceRepo:   balanceRepo,
		portfolioRepo: portfolioRepo,
		txnRepo:       txnRepo,
	}
}

func (m *MarketMaker) Run(interval time.Duration) {
	log.Printf("Market maker started (interval: %v)", interval)
	time.Sleep(3 * time.Second)

	// Look up the MARKET system user
	marketUser, err := m.userRepo.GetByUsername("MARKET")
	if err != nil {
		log.Printf("Market maker: could not find MARKET user, trades will not be recorded: %v", err)
		m.marketUserID = 0
	} else {
		m.marketUserID = marketUser.ID
		log.Printf("Market maker: using MARKET user ID %d", m.marketUserID)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	m.tick()

	for range ticker.C {
		m.tick()
	}
}

func (m *MarketMaker) tick() {
	users, err := m.userRepo.GetAll()
	if err != nil {
		log.Printf("Market maker: error fetching users: %v", err)
		return
	}

	for _, u := range users {
		// Skip the MARKET system user itself
		if u.ID == m.marketUserID {
			continue
		}

		// Apply a direct random percentage change: ±0.05% to ±0.5%
		changePct := (0.0005 + rand.Float64()*0.005) // 0.05% to 0.55%

		// Random direction
		isBuy := rand.Float64() >= 0.5
		if !isBuy {
			changePct = -changePct
		}

		// Mean-reversion bias: nudge toward baseline (10 Grub)
		if u.CurrentSharePrice > 15 && rand.Float64() < 0.35 {
			changePct = -math.Abs(changePct)
			isBuy = false
		} else if u.CurrentSharePrice < 7 && rand.Float64() < 0.35 {
			changePct = math.Abs(changePct)
			isBuy = true
		}

		newPrice := u.CurrentSharePrice * (1 + changePct)
		newPrice = math.Round(newPrice*100) / 100
		newPrice = math.Max(MinPrice, math.Min(MaxPrice, newPrice))

		if newPrice != u.CurrentSharePrice {
			if err := m.userRepo.UpdateSharePriceNoTx(u.ID, newPrice); err != nil {
				continue
			}
			if err := m.txnRepo.RecordPriceHistoryNoTx(u.ID, newPrice); err != nil {
				continue
			}

			// Record a transaction so it shows up in recent trades
			if m.marketUserID > 0 {
				priceDiff := math.Abs(newPrice - u.CurrentSharePrice)
				// Simulate a small trade: ~0.01 to 0.1 shares
				numShares := math.Round((0.01+rand.Float64()*0.09)*10000) / 10000
				totalGrub := math.Round(numShares*u.CurrentSharePrice*100) / 100
				if totalGrub < 0.01 {
					totalGrub = 0.01
				}
				_ = priceDiff

				txnType := "buy"
				if !isBuy {
					txnType = "sell"
				}

				_ = m.txnRepo.CreateNoTx(&models.Transaction{
					BuyerID:         m.marketUserID,
					StockUserID:     u.ID,
					TransactionType: txnType,
					NumShares:       numShares,
					PricePerShare:   u.CurrentSharePrice,
					TotalGrub:       totalGrub,
				})
			}
		}
	}

	m.snapshotPortfolios()
}

func (m *MarketMaker) snapshotPortfolios() {
	users, err := m.userRepo.GetAll()
	if err != nil {
		return
	}

	for _, u := range users {
		// Skip the MARKET system user
		if u.ID == m.marketUserID {
			continue
		}

		balance, err := m.balanceRepo.GetByUserID(u.ID)
		if err != nil {
			continue
		}

		holdings, err := m.portfolioRepo.GetByOwner(u.ID)
		if err != nil {
			continue
		}

		var holdingsValue float64
		for _, h := range holdings {
			stockUser, err := m.userRepo.GetByID(h.StockUserID)
			if err != nil {
				continue
			}
			holdingsValue += h.NumShares * stockUser.CurrentSharePrice
		}

		totalValue := math.Round((balance.GrubBalance+holdingsValue)*100) / 100
		_ = m.userRepo.SavePortfolioSnapshot(u.ID, totalValue, balance.GrubBalance)
	}
}
