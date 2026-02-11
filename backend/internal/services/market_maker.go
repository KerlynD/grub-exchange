package services

import (
	"database/sql"
	"grub-exchange/internal/repository"
	"log"
	"math"
	"math/rand"
	"time"
)

type MarketMaker struct {
	db            *sql.DB
	userRepo      *repository.UserRepo
	balanceRepo   *repository.BalanceRepo
	portfolioRepo *repository.PortfolioRepo
	txnRepo       *repository.TransactionRepo
	postRepo      *repository.PostRepo
	marketUserID  int
}

func NewMarketMaker(
	db *sql.DB,
	userRepo *repository.UserRepo,
	balanceRepo *repository.BalanceRepo,
	portfolioRepo *repository.PortfolioRepo,
	txnRepo *repository.TransactionRepo,
	postRepo *repository.PostRepo,
) *MarketMaker {
	return &MarketMaker{
		db:            db,
		userRepo:      userRepo,
		balanceRepo:   balanceRepo,
		portfolioRepo: portfolioRepo,
		txnRepo:       txnRepo,
		postRepo:      postRepo,
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

	// Batch-load all sentiments so we don't query per stock
	sentiments, err := m.postRepo.GetAllSentiments()
	if err != nil {
		sentiments = make(map[int]int)
	}

	for _, u := range users {
		// Skip the MARKET system user itself
		if u.ID == m.marketUserID {
			continue
		}

		// Apply a direct random percentage change: ±0.1% to ±1.0%
		changePct := (0.001 + rand.Float64()*0.01) // 0.1% to 1.1%

		// Base buy probability: 65% buy, 35% sell (bullish bias)
		buyProb := 0.65

		// Adjust buy probability based on news sentiment
		// Each net like point shifts buy probability by 2%, capped at [0.20, 0.90]
		if net, ok := sentiments[u.ID]; ok && net != 0 {
			buyProb += float64(net) * 0.02
			buyProb = math.Max(0.20, math.Min(0.90, buyProb))
		}

		isBuy := rand.Float64() < buyProb
		if !isBuy {
			changePct = -changePct
		}

		// Mean-reversion bias: nudge toward baseline (10 Grub)
		if u.CurrentSharePrice > 15 && rand.Float64() < 0.35 {
			changePct = -math.Abs(changePct)
		} else if u.CurrentSharePrice < 7 && rand.Float64() < 0.35 {
			changePct = math.Abs(changePct)
		}

		newPrice := u.CurrentSharePrice * (1 + changePct)
		newPrice = math.Round(newPrice*100) / 100
		newPrice = math.Max(MinPrice, math.Min(MaxPrice, newPrice))

		if newPrice != u.CurrentSharePrice {
			// Use a DB transaction to atomically update price + record history
			tx, err := m.db.Begin()
			if err != nil {
				log.Printf("Market maker: tx begin error for user %d: %v", u.ID, err)
				continue
			}

			_, err = tx.Exec(`UPDATE users SET current_share_price = $1 WHERE id = $2`, newPrice, u.ID)
			if err != nil {
				tx.Rollback()
				continue
			}

			_, err = tx.Exec(`INSERT INTO price_history (user_id, price, timestamp) VALUES ($1, $2, $3)`, u.ID, newPrice, time.Now())
			if err != nil {
				tx.Rollback()
				continue
			}

			if err := tx.Commit(); err != nil {
				log.Printf("Market maker: tx commit error for user %d: %v", u.ID, err)
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
